package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maynguyen24/sever/internal/models"
)

func TestTokenRepository_Operations(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewTokenRepository(db)
	now := time.Now()
	token := &models.Token{ID: 1, UserID: 42, TokenString: "refresh", ExpiresAt: now}

	mock.ExpectExec(quotedSQL(`
		INSERT INTO tokens (id, user_id, token_string, expires_at)
		VALUES ($1, $2, $3, $4)
	`)).
		WithArgs(token.ID, token.UserID, token.TokenString, token.ExpiresAt).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.StoreRefreshToken(context.Background(), token); err != nil {
		t.Fatalf("StoreRefreshToken() error = %v", err)
	}

	mock.ExpectExec(quotedSQL(`DELETE FROM tokens WHERE token_string = $1`)).
		WithArgs(token.TokenString).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.RevokeToken(context.Background(), token.TokenString); err != nil {
		t.Fatalf("RevokeToken() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, token_string, expires_at, created_at FROM tokens WHERE token_string = $1 LIMIT 1`)).
		WithArgs(token.TokenString).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token_string", "expires_at", "created_at"}).
			AddRow(token.ID, token.UserID, token.TokenString, token.ExpiresAt, now))
	got, err := repo.GetToken(context.Background(), token.TokenString)
	if err != nil || got == nil || got.UserID != 42 {
		t.Fatalf("GetToken() = %+v, %v", got, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}
