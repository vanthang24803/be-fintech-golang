package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maynguyen24/sever/configs"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
	jwtUtil "github.com/maynguyen24/sever/pkg/jwt"
	"github.com/maynguyen24/sever/pkg/snowflake"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOAuth "google.golang.org/api/oauth2/v2"
)

// Define required interfaces where used
type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, token *models.Token) error
	RevokeToken(ctx context.Context, tokenString string) error
	GetToken(ctx context.Context, tokenString string) (*models.Token, error)
}

type AuthService struct {
	userRepo   UserRepository
	tokenRepo  TokenRepository
	cfg        *configs.Config
	oauthState string
	googleCfg  *oauth2.Config
}

func NewAuthService(userRepo UserRepository, tokenRepo TokenRepository, cfg *configs.Config) *AuthService {
	googleCfg := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		cfg:        cfg,
		oauthState: "pseudo-random-state", // In production, this should be dynamic and stored in session/cookie
		googleCfg:  googleCfg,
	}
}

func (s *AuthService) GetGoogleAuthURL(ctx context.Context) string {
	return s.googleCfg.AuthCodeURL(s.oauthState, oauth2.AccessTypeOffline)
}

func (s *AuthService) HandleGoogleCallback(ctx context.Context, code string) (*models.LoginResponse, error) {
	// 1. Exchange code for token
	token, err := s.googleCfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("google code exchange failed: %w", err)
	}

	// 2. Fetch User info from Google
	client := s.googleCfg.Client(ctx, token)
	oauth2Service, err := googleOAuth.New(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 service: %w", err)
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	// 3. Find or Create User
	// Check by Google ID first
	user, err := s.userRepo.GetUserByGoogleID(ctx, userInfo.Id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// Check by Email
		user, err = s.userRepo.GetUserByEmail(ctx, userInfo.Email)
		if err != nil {
			return nil, err
		}

		if user != nil {
			// Link existing account
			if err := s.userRepo.LinkGoogleAccount(ctx, user.ID, userInfo.Id); err != nil {
				return nil, err
			}
			googleID := userInfo.Id
			user.GoogleID = &googleID
		} else {
			// Create new user
			googleID := userInfo.Id
			user = &models.User{
				ID:       snowflake.GenerateID(),
				Username: userInfo.Email, // Default username to email
				Email:    userInfo.Email,
				GoogleID: &googleID,
			}
			if err := s.userRepo.CreateUser(ctx, user); err != nil {
				return nil, err
			}

			// Update Profile with Google data
			fullName := userInfo.Name
			avatar := userInfo.Picture
			_, _ = s.userRepo.UpdateProfile(ctx, user.ID, &models.UpdateProfileRequest{
				FullName:  &fullName,
				AvatarURL: &avatar,
			})
		}
	}

	// 4. Generate Token Pair
	accessToken, refreshToken, err := jwtUtil.GenerateTokenPair(user.ID, false, s.cfg)
	if err != nil {
		return nil, err
	}

	// 5. Store Refresh Token
	tokenRecord := &models.Token{
		ID:          snowflake.GenerateID(),
		UserID:      user.ID,
		TokenString: refreshToken,
		ExpiresAt:   time.Now().Add(30 * 24 * time.Hour),
	}
	if err := s.tokenRepo.StoreRefreshToken(ctx, tokenRecord); err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		User: user,
		Tokens: &models.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// 1. Find User by email or username
	user, err := s.userRepo.GetUserByEmailOrUsername(ctx, req.Identifier, req.Identifier)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid credentials", apperr.ErrUnauthorized)
	}

	// 2. Compare Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("%w: invalid credentials", apperr.ErrUnauthorized)
	}

	// 3. Generate Token Pair (FIDO not verified by default on login)
	accessToken, refreshToken, err := jwtUtil.GenerateTokenPair(user.ID, false, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("could not generate tokens: %w", err)
	}

	// 4. Store Refresh Token
	tokenRecord := &models.Token{
		ID:          snowflake.GenerateID(),
		UserID:      user.ID,
		TokenString: refreshToken,
		ExpiresAt:   time.Now().Add(30 * 24 * time.Hour),
	}
	
	if err := s.tokenRepo.StoreRefreshToken(ctx, tokenRecord); err != nil {
		return nil, fmt.Errorf("could not save session: %w", err)
	}

	// 5. Build Response
	return &models.LoginResponse{
		User: user,
		Tokens: &models.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.TokenPair, error) {
	// 1. Verify Refresh Token Signature using jwt library directly
	token, err := jwt.ParseWithClaims(req.RefreshToken, &jwtUtil.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: unexpected signing method", apperr.ErrUnauthorized)
		}
		return []byte(s.cfg.JWTRefreshSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("%w: invalid or expired refresh token", apperr.ErrUnauthorized)
	}

	claims, ok := token.Claims.(*jwtUtil.TokenClaims)
	if !ok {
		return nil, fmt.Errorf("%w: failed to parse claims", apperr.ErrUnauthorized)
	}

	// 2. Look up the token in the Database (Anti-replay attack & session revocation check)
	_, err = s.tokenRepo.GetToken(ctx, req.RefreshToken)
	if err != nil {
		// Token is either revoked, already used, or fake
		return nil, fmt.Errorf("%w: session expired or revoked", apperr.ErrUnauthorized)
	}

	// 3. Issue a new token pair (Preserve FIDO status from old token claims)
	newAccess, newRefresh, err := jwtUtil.GenerateTokenPair(claims.UserID, claims.FIDOVerified, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// 4. Revoke Old Token (Transaction-like behavior: delete then insert)
	if err := s.tokenRepo.RevokeToken(ctx, req.RefreshToken); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	// 5. Store New Refresh Token
	tokenRecord := &models.Token{
		ID:          snowflake.GenerateID(),
		UserID:      claims.UserID,
		TokenString: newRefresh,
		ExpiresAt:   time.Now().Add(30 * 24 * time.Hour),
	}
	if err := s.tokenRepo.StoreRefreshToken(ctx, tokenRecord); err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return &models.TokenPair{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *models.LogoutRequest) error {
	if err := s.tokenRepo.RevokeToken(ctx, req.RefreshToken); err != nil {
		return fmt.Errorf("failed to revoke token during logout: %w", err)
	}
	return nil
}
