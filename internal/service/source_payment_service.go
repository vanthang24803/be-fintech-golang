package service

import (
	"context"
	"fmt"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

// SourcePaymentRepository defines the DB contract for this service
type SourcePaymentRepository interface {
	Create(ctx context.Context, source *models.SourcePayment) error
	GetAllByUserID(ctx context.Context, userID int64) ([]*models.SourcePayment, error)
	GetByID(ctx context.Context, id, userID int64) (*models.SourcePayment, error)
	Update(ctx context.Context, source *models.SourcePayment) error
	Delete(ctx context.Context, id, userID int64) error
}

type SourcePaymentService struct {
	repo SourcePaymentRepository
}

func NewSourcePaymentService(repo SourcePaymentRepository) *SourcePaymentService {
	return &SourcePaymentService{repo: repo}
}

func (s *SourcePaymentService) Create(ctx context.Context, userID int64, req *models.CreateSourcePaymentRequest) (*models.SourcePayment, error) {
	if req.Name == "" || req.Type == "" {
		return nil, fmt.Errorf("%w: Name and type are required", apperr.ErrInvalidInput)
	}

	currency := req.Currency
	if currency == "" {
		currency = "VND"
	}

	source := &models.SourcePayment{
		UserID:   userID,
		Name:     req.Name,
		Type:     req.Type,
		Balance:  req.Balance,
		Currency: currency,
	}

	if err := s.repo.Create(ctx, source); err != nil {
		return nil, fmt.Errorf("failed to create source payment: %w", err)
	}
	return source, nil
}

func (s *SourcePaymentService) GetAll(ctx context.Context, userID int64) ([]*models.SourcePayment, error) {
	sources, err := s.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source payments: %w", err)
	}
	return sources, nil
}

func (s *SourcePaymentService) GetByID(ctx context.Context, id, userID int64) (*models.SourcePayment, error) {
	source, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source payment: %w", err)
	}
	if source == nil {
		return nil, fmt.Errorf("%w: Source payment not found", apperr.ErrNotFound)
	}
	return source, nil
}

func (s *SourcePaymentService) Update(ctx context.Context, id, userID int64, req *models.UpdateSourcePaymentRequest) (*models.SourcePayment, error) {
	source, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source payment: %w", err)
	}
	if source == nil {
		return nil, fmt.Errorf("%w: Source payment not found", apperr.ErrNotFound)
	}

	source.Name = req.Name
	source.Type = req.Type
	if req.Currency != "" {
		source.Currency = req.Currency
	}

	if err := s.repo.Update(ctx, source); err != nil {
		return nil, err
	}
	return source, nil
}

func (s *SourcePaymentService) Delete(ctx context.Context, id, userID int64) error {
	source, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch source payment: %w", err)
	}
	if source == nil {
		return fmt.Errorf("%w: Source payment not found", apperr.ErrNotFound)
	}

	if err := s.repo.Delete(ctx, id, userID); err != nil {
		return err
	}
	return nil
}
