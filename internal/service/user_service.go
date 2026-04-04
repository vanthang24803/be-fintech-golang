package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/snowflake"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository interface is defined where it is used (service layer)
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmailOrUsername(email, username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByGoogleID(googleID string) (*models.User, error)
	GetUserByID(userID int64) (*models.User, error)
	GetProfileByUserID(userID int64) (*models.Profile, error)
	UpdateProfile(userID int64, req *models.UpdateProfileRequest) (*models.Profile, error)
	LinkGoogleAccount(userID int64, googleID string) error
}

// UserService struct is exported directly
type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(req *models.RegisterRequest) (*models.User, error) {
	// Simple Validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Username, Email and Password are required")
	}

	// Check if user exists
	existing, err := s.repo.GetUserByEmailOrUsername(req.Email, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "Username or Email already exists")
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
	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetProfile(userID int64) (*models.ProfileResponse, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if user == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	profile, err := s.repo.GetProfileByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}
	if profile == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Profile not found")
	}

	return &models.ProfileResponse{
		User:    user,
		Profile: profile,
	}, nil
}

func (s *UserService) UpdateProfile(userID int64, req *models.UpdateProfileRequest) (*models.Profile, error) {
	profile, err := s.repo.UpdateProfile(userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}
	if profile == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Profile not found")
	}
	return profile, nil
}
