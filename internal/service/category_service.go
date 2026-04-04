package service

import (
	"context"
	"fmt"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

// CategoryRepository defines the DB contract for this service
type CategoryRepository interface {
	Create(ctx context.Context, cat *models.Category) error
	GetAllByUserID(ctx context.Context, userID int64) ([]*models.Category, error)
	GetByID(ctx context.Context, id, userID int64) (*models.Category, error)
	GetOwnedByID(ctx context.Context, id, userID int64) (*models.Category, error)
	Update(ctx context.Context, cat *models.Category) error
	Delete(ctx context.Context, id, userID int64) error
}

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, userID int64, req *models.CreateCategoryRequest) (*models.Category, error) {
	if req.Type != models.TransactionTypeIncome && req.Type != models.TransactionTypeExpense {
		return nil, fmt.Errorf("%w: type must be 'income' or 'expense'", apperr.ErrInvalidInput)
	}

	cat := &models.Category{
		UserID: &userID,
		Name:   req.Name,
		Type:   req.Type,
		Icon:   req.Icon,
	}

	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return cat, nil
}

func (s *CategoryService) GetAll(ctx context.Context, userID int64) ([]*models.Category, error) {
	cats, err := s.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	return cats, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id, userID int64) (*models.Category, error) {
	cat, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category: %w", err)
	}
	if cat == nil {
		return nil, apperr.ErrNotFound
	}
	return cat, nil
}

func (s *CategoryService) Update(ctx context.Context, id, userID int64, req *models.UpdateCategoryRequest) (*models.Category, error) {
	if req.Type != models.TransactionTypeIncome && req.Type != models.TransactionTypeExpense {
		return nil, fmt.Errorf("%w: type must be 'income' or 'expense'", apperr.ErrInvalidInput)
	}

	// Only allow editing user-owned categories, not system defaults
	cat, err := s.repo.GetOwnedByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category: %w", err)
	}
	if cat == nil {
		return nil, fmt.Errorf("%w: category not found or cannot be modified", apperr.ErrNotFound)
	}

	cat.Name = req.Name
	cat.Type = req.Type
	cat.Icon = req.Icon

	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) Delete(ctx context.Context, id, userID int64) error {
	// Only allow deleting user-owned categories
	cat, err := s.repo.GetOwnedByID(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch category: %w", err)
	}
	if cat == nil {
		return fmt.Errorf("%w: category not found or cannot be deleted", apperr.ErrNotFound)
	}

	if err := s.repo.Delete(ctx, id, userID); err != nil {
		return err
	}
	return nil
}
