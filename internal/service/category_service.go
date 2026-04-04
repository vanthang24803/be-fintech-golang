package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
)

// CategoryRepository defines the DB contract for this service
type CategoryRepository interface {
	Create(cat *models.Category) error
	GetAllByUserID(userID int64) ([]*models.Category, error)
	GetByID(id, userID int64) (*models.Category, error)
	GetOwnedByID(id, userID int64) (*models.Category, error)
	Update(cat *models.Category) error
	Delete(id, userID int64) error
}

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(userID int64, req *models.CreateCategoryRequest) (*models.Category, error) {
	if req.Type != models.TransactionTypeIncome && req.Type != models.TransactionTypeExpense {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Type must be 'income' or 'expense'")
	}

	cat := &models.Category{
		UserID: &userID,
		Name:   req.Name,
		Type:   req.Type,
		Icon:   req.Icon,
	}

	if err := s.repo.Create(cat); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return cat, nil
}

func (s *CategoryService) GetAll(userID int64) ([]*models.Category, error) {
	cats, err := s.repo.GetAllByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	return cats, nil
}

func (s *CategoryService) GetByID(id, userID int64) (*models.Category, error) {
	cat, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category: %w", err)
	}
	if cat == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Category not found")
	}
	return cat, nil
}

func (s *CategoryService) Update(id, userID int64, req *models.UpdateCategoryRequest) (*models.Category, error) {
	if req.Type != models.TransactionTypeIncome && req.Type != models.TransactionTypeExpense {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Type must be 'income' or 'expense'")
	}

	// Only allow editing user-owned categories, not system defaults
	cat, err := s.repo.GetOwnedByID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category: %w", err)
	}
	if cat == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Category not found or cannot be modified")
	}

	cat.Name = req.Name
	cat.Type = req.Type
	cat.Icon = req.Icon

	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) Delete(id, userID int64) error {
	// Only allow deleting user-owned categories
	cat, err := s.repo.GetOwnedByID(id, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch category: %w", err)
	}
	if cat == nil {
		return fiber.NewError(fiber.StatusNotFound, "Category not found or cannot be deleted")
	}

	if err := s.repo.Delete(id, userID); err != nil {
		return err
	}
	return nil
}
