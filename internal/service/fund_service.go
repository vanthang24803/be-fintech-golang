package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

// FundRepository defines the DB contract for this service
type FundRepository interface {
	Create(ctx context.Context, fund *models.Fund) error
	GetAllByUserID(ctx context.Context, userID int64) ([]*models.Fund, error)
	GetByID(ctx context.Context, id, userID int64) (*models.Fund, error)
	Update(ctx context.Context, fund *models.Fund) error
	Delete(ctx context.Context, id, userID int64) error
	Deposit(ctx context.Context, id, userID int64, amount float64) (*models.Fund, error)
	Withdraw(ctx context.Context, id, userID int64, amount float64) (*models.Fund, error)
}

type FundService struct {
	repo FundRepository
}

func NewFundService(repo FundRepository) *FundService {
	return &FundService{repo: repo}
}

func (s *FundService) Create(ctx context.Context, userID int64, req *models.CreateFundRequest) (*models.Fund, error) {
	if req.Balance < 0 {
		return nil, fmt.Errorf("%w: Initial balance cannot be negative", apperr.ErrInvalidInput)
	}
	if req.TargetAmount < 0 {
		return nil, fmt.Errorf("%w: Target amount cannot be negative", apperr.ErrInvalidInput)
	}

	fund := &models.Fund{
		UserID:       userID,
		Name:         req.Name,
		Description:  req.Description,
		TargetAmount: req.TargetAmount,
		Balance:      req.Balance,
		Currency:     req.Currency,
	}

	if err := s.repo.Create(ctx, fund); err != nil {
		return nil, fmt.Errorf("failed to create fund: %w", err)
	}
	return fund, nil
}

func (s *FundService) GetAll(ctx context.Context, userID int64) ([]*models.Fund, error) {
	funds, err := s.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch funds: %w", err)
	}
	return funds, nil
}

func (s *FundService) GetByID(ctx context.Context, id, userID int64) (*models.Fund, error) {
	fund, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fund: %w", err)
	}
	if fund == nil {
		return nil, fmt.Errorf("%w: Fund not found", apperr.ErrNotFound)
	}
	return fund, nil
}

func (s *FundService) Update(ctx context.Context, id, userID int64, req *models.UpdateFundRequest) (*models.Fund, error) {
	if req.TargetAmount < 0 {
		return nil, fmt.Errorf("%w: Target amount cannot be negative", apperr.ErrInvalidInput)
	}

	fund, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fund: %w", err)
	}
	if fund == nil {
		return nil, fmt.Errorf("%w: Fund not found", apperr.ErrNotFound)
	}

	fund.Name = req.Name
	fund.Description = req.Description
	fund.TargetAmount = req.TargetAmount
	if req.Currency != "" {
		fund.Currency = req.Currency
	}

	if err := s.repo.Update(ctx, fund); err != nil {
		return nil, err
	}
	return fund, nil
}

func (s *FundService) Delete(ctx context.Context, id, userID int64) error {
	if err := s.repo.Delete(ctx, id, userID); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fmt.Errorf("%w: Fund not found", apperr.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s *FundService) Deposit(ctx context.Context, id, userID int64, req *models.FundTransactionRequest) (*models.Fund, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("%w: Amount must be greater than 0", apperr.ErrInvalidInput)
	}

	fund, err := s.repo.Deposit(ctx, id, userID, req.Amount)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, fmt.Errorf("%w: Fund not found", apperr.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to deposit: %w", err)
	}
	return fund, nil
}

func (s *FundService) Withdraw(ctx context.Context, id, userID int64, req *models.FundTransactionRequest) (*models.Fund, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("%w: Amount must be greater than 0", apperr.ErrInvalidInput)
	}

	fund, err := s.repo.Withdraw(ctx, id, userID, req.Amount)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, fmt.Errorf("%w: Fund not found", apperr.ErrNotFound)
		}
		if errors.Is(err, apperr.ErrInsufficientBalance) {
			return nil, fmt.Errorf("%w: Insufficient fund balance", apperr.ErrInvalidInput)
		}
		return nil, fmt.Errorf("failed to withdraw: %w", err)
	}
	return fund, nil
}
