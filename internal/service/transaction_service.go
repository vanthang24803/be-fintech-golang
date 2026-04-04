package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
	"github.com/maynguyen24/sever/pkg/queue"
)

// TransactionRepository defines the DB contract for this service
type TransactionRepository interface {
	Create(ctx context.Context, tx *models.Transaction) error
	GetAllByUserID(ctx context.Context, userID int64, filter models.TransactionFilter) ([]*models.TransactionDetail, error)
	GetByID(ctx context.Context, id, userID int64) (*models.TransactionDetail, error)
	GetRawByID(ctx context.Context, id, userID int64) (*models.Transaction, error)
	Update(ctx context.Context, old *models.Transaction, updated *models.Transaction) error
	Delete(ctx context.Context, id, userID int64) error
}

// txBudgetRepo is a subset of BudgetRepository used by TransactionService
type txBudgetRepo interface {
	GetByUserID(ctx context.Context, userID int64) ([]*models.Budget, error)
	CalculateSpending(ctx context.Context, userID int64, categoryID *int64, start, end time.Time) (float64, error)
}

// txNotificationRepo is a subset of NotificationRepository used by TransactionService
type txNotificationRepo interface {
	Create(ctx context.Context, notif *models.Notification) error
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

func (s *TransactionService) Create(ctx context.Context, userID int64, req *models.CreateTransactionRequest) (*models.Transaction, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("%w: Amount must be greater than 0", apperr.ErrInvalidInput)
	}
	if req.Type != models.TransactionTypeIncome && req.Type != models.TransactionTypeExpense {
		return nil, fmt.Errorf("%w: Type must be 'income' or 'expense'", apperr.ErrInvalidInput)
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

	if err := s.repo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if tx.Type == models.TransactionTypeExpense && s.queue != nil {
		_ = s.queue.EnqueueCheckBudget(queue.BudgetCheckPayload{
			UserID:     tx.UserID,
			CategoryID: tx.CategoryID,
		})
	}

	return tx, nil
}

func (s *TransactionService) CheckBudgets(ctx context.Context, userID int64, categoryID *int64) error {
	budgets, err := s.budgetRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, b := range budgets {
		if !b.IsActive {
			continue
		}

		if b.CategoryID != nil && (categoryID == nil || *b.CategoryID != *categoryID) {
			continue
		}

		spending, err := s.budgetRepo.CalculateSpending(ctx, userID, b.CategoryID, b.StartDate, b.EndDate)
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

		_ = s.notificationRepo.Create(ctx, &models.Notification{
			UserID:   userID,
			Source:   "budget",
			SourceID: &b.ID,
			Type:     "alert",
			Title:    title,
			Body:     body,
		})
	}
	return nil
}

func (s *TransactionService) GetAll(ctx context.Context, userID int64, query map[string]string) ([]*models.TransactionDetail, error) {
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

	txs, err := s.repo.GetAllByUserID(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	return txs, nil
}

func (s *TransactionService) GetByID(ctx context.Context, id, userID int64) (*models.TransactionDetail, error) {
	tx, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}
	if tx == nil {
		return nil, fmt.Errorf("%w: Transaction not found", apperr.ErrNotFound)
	}
	return tx, nil
}

func (s *TransactionService) Update(ctx context.Context, id, userID int64, req *models.UpdateTransactionRequest) (*models.Transaction, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("%w: Amount must be greater than 0", apperr.ErrInvalidInput)
	}
	if req.Type != models.TransactionTypeIncome && req.Type != models.TransactionTypeExpense {
		return nil, fmt.Errorf("%w: Type must be 'income' or 'expense'", apperr.ErrInvalidInput)
	}

	old, err := s.repo.GetRawByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}
	if old == nil {
		return nil, fmt.Errorf("%w: Transaction not found", apperr.ErrNotFound)
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

	if err := s.repo.Update(ctx, old, updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *TransactionService) Delete(ctx context.Context, id, userID int64) error {
	if err := s.repo.Delete(ctx, id, userID); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fmt.Errorf("%w: Transaction not found", apperr.ErrNotFound)
		}
		return err
	}
	return nil
}
