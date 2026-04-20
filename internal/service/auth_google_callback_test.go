package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/oauth2"

	"github.com/maynguyen24/sever/internal/models"
)

type rewriteTransport struct {
	base   http.RoundTripper
	target *url.URL
}

func (t *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	clone.URL.Scheme = t.target.Scheme
	clone.URL.Host = t.target.Host
	clone.Host = t.target.Host
	return t.base.RoundTrip(clone)
}

func newGoogleOAuthTestContext(t *testing.T) (context.Context, *AuthService, *stubAuthUserRepo, *stubTokenRepo) {
	t.Helper()

	userRepo := &stubAuthUserRepo{}
	tokenRepo := &stubTokenRepo{}
	svc := NewAuthService(userRepo, tokenRepo, testAuthConfig())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/token":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "google-access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			})
		case strings.HasSuffix(r.URL.Path, "/userinfo"):
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":      "google-user-123",
				"email":   "alice@example.com",
				"name":    "Alice Example",
				"picture": "https://cdn.example/avatar.png",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	targetURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("url.Parse(server.URL): %v", err)
	}

	svc.googleCfg.Endpoint.TokenURL = server.URL + "/token"
	client := &http.Client{
		Transport: &rewriteTransport{
			base:   http.DefaultTransport,
			target: targetURL,
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	return ctx, svc, userRepo, tokenRepo
}

func TestAuthService_HandleGoogleCallback_LinksExistingEmailAndStoresRefreshToken(t *testing.T) {
	t.Parallel()

	ctx, svc, userRepo, tokenRepo := newGoogleOAuthTestContext(t)

	existing := &models.User{ID: 42, Email: "alice@example.com", Username: "alice"}
	var linkedGoogleID string
	var stored *models.Token

	userRepo.getUserByGoogleIDFn = func(context.Context, string) (*models.User, error) {
		return nil, nil
	}
	userRepo.getUserByEmailFn = func(context.Context, string) (*models.User, error) {
		return existing, nil
	}
	userRepo.linkGoogleAccountFn = func(_ context.Context, userID int64, googleID string) error {
		if userID != existing.ID {
			t.Fatalf("expected link for user %d, got %d", existing.ID, userID)
		}
		linkedGoogleID = googleID
		return nil
	}
	tokenRepo.storeFn = func(_ context.Context, token *models.Token) error {
		stored = token
		return nil
	}

	resp, err := svc.HandleGoogleCallback(ctx, "valid-code")
	if err != nil {
		t.Fatalf("HandleGoogleCallback() error = %v", err)
	}
	if resp == nil || resp.User == nil || resp.Tokens == nil {
		t.Fatalf("expected login response with tokens, got %+v", resp)
	}
	if linkedGoogleID != "google-user-123" {
		t.Fatalf("expected linked google id google-user-123, got %q", linkedGoogleID)
	}
	if resp.User.GoogleID == nil || *resp.User.GoogleID != "google-user-123" {
		t.Fatalf("expected linked user google id to be set, got %+v", resp.User.GoogleID)
	}
	if stored == nil || stored.UserID != existing.ID || stored.TokenString == "" {
		t.Fatalf("expected refresh token to be stored, got %+v", stored)
	}
}

func TestAuthService_HandleGoogleCallback_CreatesUserAndUpdatesProfile(t *testing.T) {
	t.Parallel()

	ctx, svc, userRepo, tokenRepo := newGoogleOAuthTestContext(t)

	var created *models.User
	var updatedUserID int64
	var updatedProfile *models.UpdateProfileRequest

	userRepo.getUserByGoogleIDFn = func(context.Context, string) (*models.User, error) {
		return nil, nil
	}
	userRepo.getUserByEmailFn = func(context.Context, string) (*models.User, error) {
		return nil, nil
	}
	userRepo.createUserFn = func(_ context.Context, user *models.User) error {
		created = user
		return nil
	}
	userRepo.updateProfileFn = func(_ context.Context, userID int64, req *models.UpdateProfileRequest) (*models.Profile, error) {
		updatedUserID = userID
		updatedProfile = req
		return &models.Profile{UserID: userID}, nil
	}
	tokenRepo.storeFn = func(_ context.Context, token *models.Token) error {
		if token == nil || token.TokenString == "" {
			t.Fatalf("expected token to be persisted, got %+v", token)
		}
		return nil
	}

	resp, err := svc.HandleGoogleCallback(ctx, "valid-code")
	if err != nil {
		t.Fatalf("HandleGoogleCallback() error = %v", err)
	}
	if created == nil {
		t.Fatal("expected new user to be created")
	}
	if created.Email != "alice@example.com" || created.Username != "alice@example.com" {
		t.Fatalf("unexpected created user: %+v", created)
	}
	if created.GoogleID == nil || *created.GoogleID != "google-user-123" {
		t.Fatalf("expected created user google id to be set, got %+v", created.GoogleID)
	}
	if updatedProfile == nil || updatedUserID != created.ID {
		t.Fatalf("expected profile update for created user, got userID=%d req=%+v", updatedUserID, updatedProfile)
	}
	if updatedProfile.FullName == nil || *updatedProfile.FullName != "Alice Example" {
		t.Fatalf("expected full name to be set from google profile, got %+v", updatedProfile.FullName)
	}
	if updatedProfile.AvatarURL == nil || *updatedProfile.AvatarURL != "https://cdn.example/avatar.png" {
		t.Fatalf("expected avatar to be set from google profile, got %+v", updatedProfile.AvatarURL)
	}
	if resp == nil || resp.User == nil || resp.User.ID != created.ID {
		t.Fatalf("expected response user to match created user, got %+v", resp)
	}
}

func TestAuthService_HandleGoogleCallback_ExchangeFailure(t *testing.T) {
	t.Parallel()

	ctx, svc, _, _ := newGoogleOAuthTestContext(t)
	svc.googleCfg.Endpoint.TokenURL = "http://127.0.0.1:1/unreachable"

	_, err := svc.HandleGoogleCallback(ctx, "bad-code")
	if err == nil {
		t.Fatal("expected exchange failure")
	}
	if !strings.Contains(err.Error(), "google code exchange failed") {
		t.Fatalf("expected wrapped exchange error, got %v", err)
	}
}

func TestAuthService_HandleGoogleCallback_StoreTokenError(t *testing.T) {
	t.Parallel()

	ctx, svc, userRepo, tokenRepo := newGoogleOAuthTestContext(t)
	storeErr := errors.New("store failed")

	userRepo.getUserByGoogleIDFn = func(context.Context, string) (*models.User, error) {
		return &models.User{ID: 7, Email: "alice@example.com", Username: "alice"}, nil
	}
	tokenRepo.storeFn = func(context.Context, *models.Token) error {
		return storeErr
	}

	_, err := svc.HandleGoogleCallback(ctx, "valid-code")
	if !errors.Is(err, storeErr) {
		t.Fatalf("expected store error, got %v", err)
	}
}
