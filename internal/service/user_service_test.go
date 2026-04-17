package service

import (
	"context"
	"errors"
	"testing"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
	"golang.org/x/crypto/bcrypt"
)

type stubUserRepo struct {
	createUserFn               func(context.Context, *models.User) error
	getUserByEmailOrUsernameFn func(context.Context, string, string) (*models.User, error)
	getUserByIDFn              func(context.Context, int64) (*models.User, error)
	getProfileByUserIDFn       func(context.Context, int64) (*models.Profile, error)
	updateProfileFn            func(context.Context, int64, *models.UpdateProfileRequest) (*models.Profile, error)
}

func (s *stubUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	return s.createUserFn(ctx, user)
}

func (s *stubUserRepo) GetUserByEmailOrUsername(ctx context.Context, email, username string) (*models.User, error) {
	return s.getUserByEmailOrUsernameFn(ctx, email, username)
}

func (s *stubUserRepo) GetUserByEmail(context.Context, string) (*models.User, error) {
	return nil, nil
}

func (s *stubUserRepo) GetUserByGoogleID(context.Context, string) (*models.User, error) {
	return nil, nil
}

func (s *stubUserRepo) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	return s.getUserByIDFn(ctx, userID)
}

func (s *stubUserRepo) GetProfileByUserID(ctx context.Context, userID int64) (*models.Profile, error) {
	return s.getProfileByUserIDFn(ctx, userID)
}

func (s *stubUserRepo) UpdateProfile(ctx context.Context, userID int64, req *models.UpdateProfileRequest) (*models.Profile, error) {
	return s.updateProfileFn(ctx, userID, req)
}

func (s *stubUserRepo) LinkGoogleAccount(context.Context, int64, string) error {
	return nil
}

func TestUserService_RegisterUser(t *testing.T) {
	t.Parallel()

	repo := &stubUserRepo{}
	repo.getUserByEmailOrUsernameFn = func(context.Context, string, string) (*models.User, error) {
		return nil, nil
	}
	var created *models.User
	repo.createUserFn = func(ctx context.Context, user *models.User) error {
		created = user
		return nil
	}

	svc := NewUserService(repo)
	user, err := svc.RegisterUser(context.Background(), &models.RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("RegisterUser returned error: %v", err)
	}
	if created == nil || user == nil {
		t.Fatal("expected user to be created")
	}
	if created.Username != "alice" || created.Email != "alice@example.com" {
		t.Fatalf("unexpected user persisted: %+v", created)
	}
	if created.PasswordHash == "secret123" {
		t.Fatal("password was not hashed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(created.PasswordHash), []byte("secret123")); err != nil {
		t.Fatalf("stored password hash did not match original password: %v", err)
	}
}

func TestUserService_RegisterUser_ValidationAndConflict(t *testing.T) {
	t.Parallel()

	svc := NewUserService(&stubUserRepo{
		getUserByEmailOrUsernameFn: func(context.Context, string, string) (*models.User, error) {
			return &models.User{ID: 1}, nil
		},
		createUserFn: func(context.Context, *models.User) error {
			t.Fatal("CreateUser should not be called for invalid/conflicting register attempts")
			return nil
		},
	})

	invalidCases := []struct {
		name string
		req  *models.RegisterRequest
	}{
		{name: "missing username", req: &models.RegisterRequest{Email: "a@example.com", Password: "secret123"}},
		{name: "missing email", req: &models.RegisterRequest{Username: "alice", Password: "secret123"}},
		{name: "missing password", req: &models.RegisterRequest{Username: "alice", Email: "a@example.com"}},
	}

	for _, tt := range invalidCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := svc.RegisterUser(context.Background(), tt.req)
			if !errors.Is(err, apperr.ErrInvalidInput) {
				t.Fatalf("expected invalid input, got %v", err)
			}
		})
	}

	_, err := svc.RegisterUser(context.Background(), &models.RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "secret123",
	})
	if !errors.Is(err, apperr.ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestUserService_GetProfile(t *testing.T) {
	t.Parallel()

	fullName := "Alice"
	repo := &stubUserRepo{}
	repo.getUserByIDFn = func(context.Context, int64) (*models.User, error) {
		return &models.User{ID: 7, Username: "alice", Email: "alice@example.com"}, nil
	}
	repo.getProfileByUserIDFn = func(context.Context, int64) (*models.Profile, error) {
		return &models.Profile{ID: 9, UserID: 7, FullName: &fullName}, nil
	}

	svc := NewUserService(repo)
	resp, err := svc.GetProfile(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetProfile returned error: %v", err)
	}
	if resp.User == nil || resp.Profile == nil || resp.Profile.FullName == nil || *resp.Profile.FullName != "Alice" {
		t.Fatalf("unexpected profile response: %+v", resp)
	}

	repo.getUserByIDFn = func(context.Context, int64) (*models.User, error) {
		return nil, nil
	}
	_, err = svc.GetProfile(context.Background(), 7)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected user not found, got %v", err)
	}

	repo.getUserByIDFn = func(context.Context, int64) (*models.User, error) {
		return &models.User{ID: 7}, nil
	}
	repo.getProfileByUserIDFn = func(context.Context, int64) (*models.Profile, error) {
		return nil, nil
	}
	_, err = svc.GetProfile(context.Background(), 7)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected profile not found, got %v", err)
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	t.Parallel()

	profile := &models.Profile{ID: 1, UserID: 7}
	svc := NewUserService(&stubUserRepo{
		updateProfileFn: func(context.Context, int64, *models.UpdateProfileRequest) (*models.Profile, error) {
			return profile, nil
		},
	})

	got, err := svc.UpdateProfile(context.Background(), 7, &models.UpdateProfileRequest{})
	if err != nil {
		t.Fatalf("UpdateProfile returned error: %v", err)
	}
	if got != profile {
		t.Fatalf("unexpected profile returned: %+v", got)
	}

	svc = NewUserService(&stubUserRepo{
		updateProfileFn: func(context.Context, int64, *models.UpdateProfileRequest) (*models.Profile, error) {
			return nil, nil
		},
	})
	_, err = svc.UpdateProfile(context.Background(), 7, &models.UpdateProfileRequest{})
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}
