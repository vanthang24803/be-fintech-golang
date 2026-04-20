package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maynguyen24/sever/internal/models"
)

func TestCategoryRepository_Operations(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewCategoryRepository(db)
	now := time.Now()
	userID := int64(42)
	icon := "wallet"
	cat := &models.Category{UserID: &userID, Name: "Salary", Type: "income", Icon: &icon}

	mock.ExpectQuery(quotedSQL(`
		INSERT INTO categories (id, user_id, name, type, icon)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`)).
		WithArgs(sqlmock.AnyArg(), cat.UserID, cat.Name, cat.Type, cat.Icon).
		WillReturnRows(timestampRows(now))
	if err := repo.Create(context.Background(), cat); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT id, user_id, name, type, icon, created_at, updated_at
		FROM categories
		WHERE user_id = $1 OR user_id IS NULL
		ORDER BY user_id NULLS FIRST, name ASC
	`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows(categoryCols).
			AddRow(int64(1), nil, "System", "expense", nil, now, now).
			AddRow(int64(2), userID, "Salary", "income", icon, now, now))
	list, err := repo.GetAllByUserID(context.Background(), userID)
	if err != nil || len(list) != 2 {
		t.Fatalf("GetAllByUserID() = %+v, %v", list, err)
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT id, user_id, name, type, icon, created_at, updated_at
		FROM categories
		WHERE id = $1 AND (user_id = $2 OR user_id IS NULL)
		LIMIT 1
	`)).
		WithArgs(int64(2), userID).
		WillReturnRows(sqlmock.NewRows(categoryCols).
			AddRow(int64(2), userID, "Salary", "income", icon, now, now))
	got, err := repo.GetByID(context.Background(), 2, userID)
	if err != nil || got == nil || got.Name != "Salary" {
		t.Fatalf("GetByID() = %+v, %v", got, err)
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT id, user_id, name, type, icon, created_at, updated_at
		FROM categories
		WHERE id = $1 AND user_id = $2
		LIMIT 1
	`)).
		WithArgs(int64(2), userID).
		WillReturnRows(sqlmock.NewRows(categoryCols).
			AddRow(int64(2), userID, "Salary", "income", icon, now, now))
	got, err = repo.GetOwnedByID(context.Background(), 2, userID)
	if err != nil || got == nil || got.ID != 2 {
		t.Fatalf("GetOwnedByID() = %+v, %v", got, err)
	}

	mock.ExpectQuery(quotedSQL(`
		UPDATE categories
		SET name = $1, type = $2, icon = $3, updated_at = NOW()
		WHERE id = $4 AND user_id = $5
		RETURNING updated_at
	`)).
		WithArgs(cat.Name, cat.Type, cat.Icon, cat.ID, cat.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(now))
	if err := repo.Update(context.Background(), cat); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	mock.ExpectExec(quotedSQL(`DELETE FROM categories WHERE id = $1 AND user_id = $2`)).
		WithArgs(cat.ID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.Delete(context.Background(), cat.ID, userID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT id, user_id, name, type, icon, created_at, updated_at
		FROM categories
		WHERE id = $1 AND (user_id = $2 OR user_id IS NULL)
		LIMIT 1
	`)).
		WithArgs(int64(99), userID).
		WillReturnError(sql.ErrNoRows)
	got, err = repo.GetByID(context.Background(), 99, userID)
	if err != nil || got != nil {
		t.Fatalf("expected nil on missing category, got %+v err=%v", got, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}
