package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
)

// SourcePaymentRepository defines the DB contract for this service
type SourcePaymentRepository interface {
	Create(source *models.SourcePayment) error
	GetAllByUserID(userID int64) ([]*models.SourcePayment, error)
	GetByID(id, userID int64) (*models.SourcePayment, error)
	Update(source *models.SourcePayment) error
	Delete(id, userID int64) error
}

type SourcePaymentService struct {
	repo SourcePaymentRepository
}

func NewSourcePaymentService(repo SourcePaymentRepository) *SourcePaymentService {
	return &SourcePaymentService{repo: repo}
}

func (s *SourcePaymentService) Create(userID int64, req *models.CreateSourcePaymentRequest) (*models.SourcePayment, error) {
	if req.Name == "" || req.Type == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Name and type are required")
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

	if err := s.repo.Create(source); err != nil {
		return nil, fmt.Errorf("failed to create source payment: %w", err)
	}
	return source, nil
}

func (s *SourcePaymentService) GetAll(userID int64) ([]*models.SourcePayment, error) {
	sources, err := s.repo.GetAllByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source payments: %w", err)
	}
	return sources, nil
}

func (s *SourcePaymentService) GetByID(id, userID int64) (*models.SourcePayment, error) {
	source, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source payment: %w", err)
	}
	if source == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Source payment not found")
	}
	return source, nil
}

func (s *SourcePaymentService) Update(id, userID int64, req *models.UpdateSourcePaymentRequest) (*models.SourcePayment, error) {
	// Check ownership
	source, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source payment: %w", err)
	}
	if source == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Source payment not found")
	}

	source.Name = req.Name
	source.Type = req.Type
	if req.Currency != "" {
		source.Currency = req.Currency
	}

	if err := s.repo.Update(source); err != nil {
		return nil, err
	}
	return source, nil
}

func (s *SourcePaymentService) Delete(id, userID int64) error {
	// Check ownership first
	source, err := s.repo.GetByID(id, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch source payment: %w", err)
	}
	if source == nil {
		return fiber.NewError(fiber.StatusNotFound, "Source payment not found")
	}

	if err := s.repo.Delete(id, userID); err != nil {
		return err
	}
	return nil
}
