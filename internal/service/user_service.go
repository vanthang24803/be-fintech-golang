package service

import (
	"context"
	"fmt"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
	"github.com/maynguyen24/sever/pkg/snowflake"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository interface is defined where it is used (service layer)
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmailOrUsername(ctx context.Context, email, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	GetProfileByUserID(ctx context.Context, userID int64) (*models.Profile, error)
	UpdateProfile(ctx context.Context, userID int64, req *models.UpdateProfileRequest) (*models.Profile, error)
	LinkGoogleAccount(ctx context.Context, userID int64, googleID string) error
}

// UserService struct is exported directly
type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	// Simple Validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("%w: Username, Email and Password are required", apperr.ErrInvalidInput)
	}

	// Check if user exists
	existing, err := s.repo.GetUserByEmailOrUsername(ctx, req.Email, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, apperr.ErrConflict
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create entity
	user := &models.User{
		ID:           snowflake.GenerateID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
	}

	// Save to DB
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (*models.ProfileResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if user == nil {
		return nil, apperr.ErrNotFound
	}

	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}
	if profile == nil {
		return nil, apperr.ErrNotFound
	}

	return &models.ProfileResponse{
		User:    user,
		Profile: profile,
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID int64, req *models.UpdateProfileRequest) (*models.Profile, error) {
	profile, err := s.repo.UpdateProfile(ctx, userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}
	if profile == nil {
		return nil, apperr.ErrNotFound
	}
	return profile, nil
}
