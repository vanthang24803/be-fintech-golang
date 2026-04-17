package service

import (
	"context"
	"errors"
	"testing"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

type stubCategoryRepo struct {
	createFn       func(context.Context, *models.Category) error
	getByIDFn      func(context.Context, int64, int64) (*models.Category, error)
	getOwnedByIDFn func(context.Context, int64, int64) (*models.Category, error)
	updateFn       func(context.Context, *models.Category) error
	deleteFn       func(context.Context, int64, int64) error
}

func (s *stubCategoryRepo) Create(ctx context.Context, cat *models.Category) error {
	return s.createFn(ctx, cat)
}

func (s *stubCategoryRepo) GetAllByUserID(context.Context, int64) ([]*models.Category, error) {
	return nil, nil
}

func (s *stubCategoryRepo) GetByID(ctx context.Context, id, userID int64) (*models.Category, error) {
	return s.getByIDFn(ctx, id, userID)
}

func (s *stubCategoryRepo) GetOwnedByID(ctx context.Context, id, userID int64) (*models.Category, error) {
	return s.getOwnedByIDFn(ctx, id, userID)
}

func (s *stubCategoryRepo) Update(ctx context.Context, cat *models.Category) error {
	return s.updateFn(ctx, cat)
}

func (s *stubCategoryRepo) Delete(ctx context.Context, id, userID int64) error {
	return s.deleteFn(ctx, id, userID)
}

func TestCategoryService_CreateAndValidateType(t *testing.T) {
	t.Parallel()

	icon := "icon"
	userID := int64(12)
	repo := &stubCategoryRepo{
		createFn: func(ctx context.Context, cat *models.Category) error {
			if cat.UserID == nil || *cat.UserID != userID || cat.Name != "Food" || cat.Type != models.TransactionTypeExpense || cat.Icon != &icon {
				t.Fatalf("unexpected category passed to repo: %+v", cat)
			}
			return nil
		},
	}
	svc := NewCategoryService(repo)

	cat, err := svc.Create(context.Background(), userID, &models.CreateCategoryRequest{
		Name: "Food",
		Type: models.TransactionTypeExpense,
		Icon: &icon,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if cat.UserID == nil || *cat.UserID != userID || cat.Name != "Food" {
		t.Fatalf("unexpected category returned: %+v", cat)
	}

	_, err = svc.Create(context.Background(), userID, &models.CreateCategoryRequest{Name: "x", Type: "invalid"})
	if !errors.Is(err, apperr.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestCategoryService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	svc := NewCategoryService(&stubCategoryRepo{
		getByIDFn: func(context.Context, int64, int64) (*models.Category, error) {
			return nil, nil
		},
	})

	_, err := svc.GetByID(context.Background(), 1, 2)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestCategoryService_Update(t *testing.T) {
	t.Parallel()

	icon := "updated"
	owned := int64(42)
	repo := &stubCategoryRepo{
		getOwnedByIDFn: func(context.Context, int64, int64) (*models.Category, error) {
			return &models.Category{ID: 1, UserID: &owned, Name: "Old", Type: models.TransactionTypeIncome}, nil
		},
		updateFn: func(ctx context.Context, cat *models.Category) error {
			if cat.Name != "New" || cat.Type != models.TransactionTypeExpense || cat.Icon != &icon {
				t.Fatalf("unexpected category passed to update: %+v", cat)
			}
			return nil
		},
	}
	svc := NewCategoryService(repo)

	cat, err := svc.Update(context.Background(), 1, owned, &models.UpdateCategoryRequest{
		Name: "New",
		Type: models.TransactionTypeExpense,
		Icon: &icon,
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if cat.Name != "New" || cat.Type != models.TransactionTypeExpense {
		t.Fatalf("unexpected category returned: %+v", cat)
	}

	_, err = svc.Update(context.Background(), 1, owned, &models.UpdateCategoryRequest{Name: "New", Type: "invalid"})
	if !errors.Is(err, apperr.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestCategoryService_Update_NotOwned(t *testing.T) {
	t.Parallel()

	svc := NewCategoryService(&stubCategoryRepo{
		getOwnedByIDFn: func(context.Context, int64, int64) (*models.Category, error) {
			return nil, nil
		},
		updateFn: func(context.Context, *models.Category) error {
			t.Fatal("repo.Update should not be called when category is not owned")
			return nil
		},
	})

	_, err := svc.Update(context.Background(), 1, 2, &models.UpdateCategoryRequest{Name: "New", Type: models.TransactionTypeIncome})
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestCategoryService_Delete(t *testing.T) {
	t.Parallel()

	owned := int64(42)
	deleted := false
	svc := NewCategoryService(&stubCategoryRepo{
		getOwnedByIDFn: func(context.Context, int64, int64) (*models.Category, error) {
			return &models.Category{ID: 1, UserID: &owned, Name: "Food", Type: models.TransactionTypeExpense}, nil
		},
		deleteFn: func(context.Context, int64, int64) error {
			deleted = true
			return nil
		},
	})

	if err := svc.Delete(context.Background(), 1, owned); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if !deleted {
		t.Fatal("repo.Delete was not called")
	}
}

func TestCategoryService_Delete_NotOwned(t *testing.T) {
	t.Parallel()

	svc := NewCategoryService(&stubCategoryRepo{
		getOwnedByIDFn: func(context.Context, int64, int64) (*models.Category, error) {
			return nil, nil
		},
		deleteFn: func(context.Context, int64, int64) error {
			t.Fatal("repo.Delete should not be called when category is not owned")
			return nil
		},
	})

	err := svc.Delete(context.Background(), 1, 2)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}
