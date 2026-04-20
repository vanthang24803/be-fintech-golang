package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maynguyen24/sever/internal/models"
)

func TestUserRepository_GetUserByEmailOrUsername(t *testing.T) {
	t.Parallel()

	now := time.Now()
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)

	mock.ExpectQuery(quotedSQL("SELECT id, username, email, password_hash, google_id, created_at, updated_at FROM users WHERE email = $1 OR username = $2 LIMIT 1")).
		WithArgs("a@example.com", "alice").
		WillReturnRows(sqlmock.NewRows(userCols).
			AddRow(int64(1), "alice", "a@example.com", "hash", nil, now, now))
	mock.ExpectQuery(quotedSQL("SELECT id, username, email, password_hash, google_id, created_at, updated_at FROM users WHERE email = $1 OR username = $2 LIMIT 1")).
		WithArgs("missing@example.com", "missing").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery(quotedSQL("SELECT id, username, email, password_hash, google_id, created_at, updated_at FROM users WHERE email = $1 OR username = $2 LIMIT 1")).
		WithArgs("err@example.com", "err").
		WillReturnError(sql.ErrConnDone)

	user, err := repo.GetUserByEmailOrUsername(context.Background(), "a@example.com", "alice")
	if err != nil || user == nil || user.Email != "a@example.com" {
		t.Fatalf("GetUserByEmailOrUsername() = %+v, %v", user, err)
	}
	user, err = repo.GetUserByEmailOrUsername(context.Background(), "missing@example.com", "missing")
	if err != nil || user != nil {
		t.Fatalf("expected nil user on not found, got %+v err=%v", user, err)
	}
	_, err = repo.GetUserByEmailOrUsername(context.Background(), "err@example.com", "err")
	if !errors.Is(err, sql.ErrConnDone) {
		t.Fatalf("expected db error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}

func TestUserRepository_OtherGetters(t *testing.T) {
	t.Parallel()

	now := time.Now()
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)

	mock.ExpectQuery(quotedSQL("SELECT id, username, email, password_hash, google_id, created_at, updated_at FROM users WHERE email = $1 LIMIT 1")).
		WithArgs("a@example.com").
		WillReturnRows(sqlmock.NewRows(userCols).
			AddRow(int64(1), "alice", "a@example.com", "hash", nil, now, now))
	mock.ExpectQuery(quotedSQL("SELECT id, username, email, password_hash, google_id, created_at, updated_at FROM users WHERE google_id = $1 LIMIT 1")).
		WithArgs("google-1").
		WillReturnRows(sqlmock.NewRows(userCols).
			AddRow(int64(1), "alice", "a@example.com", "hash", "google-1", now, now))
	mock.ExpectQuery(quotedSQL("SELECT id, username, email, password_hash, google_id, created_at, updated_at FROM users WHERE id = $1 LIMIT 1")).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows(userCols).
			AddRow(int64(42), "alice", "a@example.com", "hash", nil, now, now))
	mock.ExpectQuery(quotedSQL("SELECT id, user_id, full_name, avatar_url, phone_number, date_of_birth, created_at, updated_at FROM profiles WHERE user_id = $1 LIMIT 1")).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows(profileCols).
			AddRow(int64(10), int64(42), "Alice", nil, nil, nil, now, now))
	mock.ExpectQuery(quotedSQL("SELECT id, user_id, full_name, avatar_url, phone_number, date_of_birth, created_at, updated_at FROM profiles WHERE user_id = $1 LIMIT 1")).
		WithArgs(int64(99)).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByEmail(context.Background(), "a@example.com")
	if err != nil || user == nil || user.Email != "a@example.com" {
		t.Fatalf("GetUserByEmail() = %+v, %v", user, err)
	}
	user, err = repo.GetUserByGoogleID(context.Background(), "google-1")
	if err != nil || user == nil || user.Username != "alice" {
		t.Fatalf("GetUserByGoogleID() = %+v, %v", user, err)
	}
	user, err = repo.GetUserByID(context.Background(), 42)
	if err != nil || user == nil || user.ID != 42 {
		t.Fatalf("GetUserByID() = %+v, %v", user, err)
	}
	profile, err := repo.GetProfileByUserID(context.Background(), 42)
	if err != nil || profile == nil || profile.UserID != 42 {
		t.Fatalf("GetProfileByUserID() = %+v, %v", profile, err)
	}
	profile, err = repo.GetProfileByUserID(context.Background(), 99)
	if err != nil || profile != nil {
		t.Fatalf("expected nil profile on not found, got %+v err=%v", profile, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}

func TestUserRepository_CreateUserAndUpdateProfile(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	now := time.Now()
	user := &models.User{ID: 42, Username: "alice", Email: "a@example.com", PasswordHash: "hash"}

	mock.ExpectBegin()
	mock.ExpectQuery(quotedSQL(`
		INSERT INTO users (id, username, email, password_hash, google_id) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING created_at, updated_at
	`)).
		WithArgs(user.ID, user.Username, user.Email, user.PasswordHash, user.GoogleID).
		WillReturnRows(timestampRows(now))
	mock.ExpectExec(quotedSQL(`
		INSERT INTO profiles (id, user_id) 
		VALUES ($1, $2)
	`)).
		WithArgs(sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	if user.CreatedAt.IsZero() || user.UpdatedAt.IsZero() {
		t.Fatalf("expected timestamps to be set: %+v", user)
	}

	fullName := "Alice"
	mock.ExpectQuery(quotedSQL(`
		UPDATE profiles
		SET
			full_name    = COALESCE($1, full_name),
			avatar_url   = COALESCE($2, avatar_url),
			phone_number = COALESCE($3, phone_number),
			date_of_birth = COALESCE($4, date_of_birth),
			updated_at   = NOW()
		WHERE user_id = $5
		RETURNING id, user_id, full_name, avatar_url, phone_number, date_of_birth, created_at, updated_at
	`)).
		WithArgs(&fullName, (*string)(nil), (*string)(nil), (*time.Time)(nil), int64(42)).
		WillReturnRows(sqlmock.NewRows(profileCols).
			AddRow(int64(1), int64(42), fullName, nil, nil, nil, now, now))

	profile, err := repo.UpdateProfile(context.Background(), 42, &models.UpdateProfileRequest{FullName: &fullName})
	if err != nil || profile == nil || profile.UserID != 42 {
		t.Fatalf("UpdateProfile() = %+v, %v", profile, err)
	}

	mock.ExpectExec(quotedSQL("UPDATE users SET google_id = $1, updated_at = NOW() WHERE id = $2")).
		WithArgs("google-1", int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.LinkGoogleAccount(context.Background(), 42, "google-1"); err != nil {
		t.Fatalf("LinkGoogleAccount() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}
