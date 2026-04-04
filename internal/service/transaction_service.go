package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/queue"
)

// TransactionRepository defines the DB contract for this service
type TransactionRepository interface {
	Create(tx *models.Transaction) error
	GetAllByUserID(userID int64, filter models.TransactionFilter) ([]*models.TransactionDetail, error)
	GetByID(id, userID int64) (*models.TransactionDetail, error)
	GetRawByID(id, userID int64) (*models.Transaction, error)
	Update(old *models.Transaction, updated *models.Transaction) error
	Delete(id, userID int64) error
}

// txBudgetRepo is a subset of BudgetRepository used by TransactionService
type txBudgetRepo interface {
	GetByUserID(userID int64) ([]*models.Budget, error)
	CalculateSpending(userID int64, categoryID *int64, start, end time.Time) (float64, error)
}

// txNotificationRepo is a subset of NotificationRepository used by TransactionService
type txNotificationRepo interface {
	Create(notif *models.Notification) error
}

type TransactionService struct {
	repo             TransactionRepository
	budgetRepo       txBudgetRepo
	notificationRepo txNotificationRepo
	queue            *queue.Client // Add queue client
}

func NewTransactionService(
	repo TransactionRepository,
	budgetRepo txBudgetRepo,
	notificationRepo txNotificationRepo,
	queue *queue.Client,
) *TransactionService {
	return &TransactionService{
		repo:             repo,
		budgetRepo:       budgetRepo,
		notificationRepo: notificationRepo,
		queue:            queue,
	}
}

func (s *TransactionService) Create(userID int64, req *models.CreateTransactionRequest) (*models.Transaction, error) {
	if req.Amount <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Amount must be greater than 0")
	}
	if req.Type != "income" && req.Type != "expense" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Type must be 'income' or 'expense'")
	}

	tx := &models.Transaction{
		UserID:          userID,
		SourcePaymentID: req.SourcePaymentID,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		Type:            req.Type,
		Description:     req.Description,
		TransactionDate: req.TransactionDate,
	}

	if err := s.repo.Create(tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Internal budget check (Enqueue Task)
	if tx.Type == "expense" && s.queue != nil {
		_ = s.queue.EnqueueCheckBudget(queue.BudgetCheckPayload{
			UserID:     tx.UserID,
			CategoryID: tx.CategoryID,
		})
	}

	return tx, nil
}

// CheckBudgets verifies thresholds and sends notifications (Called by Worker)
func (s *TransactionService) CheckBudgets(userID int64, categoryID *int64) {
	budgets, err := s.budgetRepo.GetByUserID(userID)
	if err != nil {
		return
	}

	for _, b := range budgets {
		if !b.IsActive {
			continue
		}

		// Check if budget applies (category match or global budget)
		if b.CategoryID != nil && (categoryID == nil || *b.CategoryID != *categoryID) {
			continue
		}

		spending, err := s.budgetRepo.CalculateSpending(userID, b.CategoryID, b.StartDate, b.EndDate)
		if err != nil {
			continue
		}

		progress := (spending / b.Amount) * 100
		var title, body string

		if progress >= 100 {
			title = "Budget Exceeded! ⚠️"
			body = fmt.Sprintf("You have spent %.2f, which is over your budget of %.2f.", spending, b.Amount)
		} else if progress >= 80 {
			title = "Budget Warning 📉"
			body = fmt.Sprintf("You have spent %.2f, reaching %.0f%% of your budget (%.2f).", spending, progress, b.Amount)
		} else {
			continue // No notification needed
		}

		// Create notification
		_ = s.notificationRepo.Create(&models.Notification{
			UserID:   userID,
			Source:   "budget",
			SourceID: &b.ID,
			Type:     "alert",
			Title:    title,
			Body:     body,
		})
	}
}

func (s *TransactionService) GetAll(userID int64, query map[string]string) ([]*models.TransactionDetail, error) {
	filter := models.TransactionFilter{}

	if t := query["type"]; t != "" {
		filter.Type = t
	}
	if cid := query["category_id"]; cid != "" {
		if id, err := strconv.ParseInt(cid, 10, 64); err == nil {
			filter.CategoryID = id
		}
	}
	if sid := query["source_id"]; sid != "" {
		if id, err := strconv.ParseInt(sid, 10, 64); err == nil {
			filter.SourcePaymentID = id
		}
	}

	txs, err := s.repo.GetAllByUserID(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	return txs, nil
}

func (s *TransactionService) GetByID(id, userID int64) (*models.TransactionDetail, error) {
	tx, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}
	if tx == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Transaction not found")
	}
	return tx, nil
}

func (s *TransactionService) Update(id, userID int64, req *models.UpdateTransactionRequest) (*models.Transaction, error) {
	if req.Amount <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Amount must be greater than 0")
	}
	if req.Type != "income" && req.Type != "expense" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Type must be 'income' or 'expense'")
	}

	old, err := s.repo.GetRawByID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}
	if old == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Transaction not found")
	}

	updated := &models.Transaction{
		ID:              id,
		UserID:          userID,
		SourcePaymentID: req.SourcePaymentID,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		Type:            req.Type,
		Description:     req.Description,
		TransactionDate: req.TransactionDate,
	}

	if err := s.repo.Update(old, updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *TransactionService) Delete(id, userID int64) error {
	if err := s.repo.Delete(id, userID); err != nil {
		if err.Error() == "transaction not found" {
			return fiber.NewError(fiber.StatusNotFound, "Transaction not found")
		}
		return err
	}
	return nil
}
