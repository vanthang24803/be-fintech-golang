package service

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
)

// BudgetRepository defines the DB contract for this service
type BudgetRepository interface {
	Create(budget *models.Budget) error
	GetByUserID(userID int64) ([]*models.Budget, error)
	GetByID(id, userID int64) (*models.Budget, error)
	Update(id, userID int64, req *models.UpdateBudgetRequest) (*models.Budget, error)
	Delete(id, userID int64) error
	CalculateSpending(userID int64, categoryID *int64, start, end time.Time) (float64, error)
}

type BudgetService struct {
	repo BudgetRepository
}

func NewBudgetService(repo BudgetRepository) *BudgetService {
	return &BudgetService{repo: repo}
}

// Create sets up a new budget and determines start/end dates
func (s *BudgetService) Create(userID int64, req *models.CreateBudgetRequest) (*models.Budget, error) {
	if req.Amount <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Budget amount must be greater than zero")
	}

	var start, end time.Time
	now := time.Now()

	switch req.Period {
	case "monthly":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, -1) // Corrected: last day of month
		// Ensure end date is until the very end of the day (23:59:59)
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, now.Location())
	case "weekly":
		// Assume week starts on Monday
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday
		}
		start = now.AddDate(0, 0, -(weekday - 1))
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 6)
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, now.Location())
	default:
		return nil, fiber.NewError(fiber.StatusBadRequest, "Unsupported budget period (only 'monthly' or 'weekly' supported)")
	}

	budget := &models.Budget{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Period:     req.Period,
		StartDate:  start,
		EndDate:    end,
		IsActive:   true,
	}

	if err := s.repo.Create(budget); err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	return budget, nil
}

func (s *BudgetService) GetList(userID int64) ([]*models.BudgetResponse, error) {
	budgets, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch budgets: %w", err)
	}

	var responses []*models.BudgetResponse
	for _, b := range budgets {
		spending, err := s.repo.CalculateSpending(userID, b.CategoryID, b.StartDate, b.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate spending for budget %d: %w", b.ID, err)
		}

		remaining := b.Amount - spending
		if remaining < 0 {
			remaining = 0
		}

		progress := (spending / b.Amount) * 100
		if progress > 100 {
			progress = 100
		}

		responses = append(responses, &models.BudgetResponse{
			Budget:          *b,
			CurrentSpending: spending,
			RemainingAmount: remaining,
			ProgressPercent: progress,
		})
	}

	return responses, nil
}

func (s *BudgetService) GetDetail(id, userID int64) (*models.BudgetResponse, error) {
	b, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch budget: %w", err)
	}
	if b == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Budget not found")
	}

	spending, err := s.repo.CalculateSpending(userID, b.CategoryID, b.StartDate, b.EndDate)
	if err != nil {
		return nil, err
	}

	return &models.BudgetResponse{
		Budget:          *b,
		CurrentSpending: spending,
		RemainingAmount: b.Amount - spending,
		ProgressPercent: (spending / b.Amount) * 100,
	}, nil
}

func (s *BudgetService) Update(id, userID int64, req *models.UpdateBudgetRequest) (*models.Budget, error) {
	budget, err := s.repo.Update(id, userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}
	if budget == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Budget not found")
	}
	return budget, nil
}

func (s *BudgetService) Delete(id, userID int64) error {
	if err := s.repo.Delete(id, userID); err != nil {
		if err.Error() == "budget not found" {
			return fiber.NewError(fiber.StatusNotFound, "Budget not found")
		}
		return err
	}
	return nil
}
