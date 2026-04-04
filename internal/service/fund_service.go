package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
)

// FundRepository defines the DB contract for this service
type FundRepository interface {
	Create(fund *models.Fund) error
	GetAllByUserID(userID int64) ([]*models.Fund, error)
	GetByID(id, userID int64) (*models.Fund, error)
	Update(fund *models.Fund) error
	Delete(id, userID int64) error
	Deposit(id, userID int64, amount float64) (*models.Fund, error)
	Withdraw(id, userID int64, amount float64) (*models.Fund, error)
}

type FundService struct {
	repo FundRepository
}

func NewFundService(repo FundRepository) *FundService {
	return &FundService{repo: repo}
}

func (s *FundService) Create(userID int64, req *models.CreateFundRequest) (*models.Fund, error) {
	if req.Balance < 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Initial balance cannot be negative")
	}
	if req.TargetAmount < 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Target amount cannot be negative")
	}

	fund := &models.Fund{
		UserID:       userID,
		Name:         req.Name,
		Description:  req.Description,
		TargetAmount: req.TargetAmount,
		Balance:      req.Balance,
		Currency:     req.Currency,
	}

	if err := s.repo.Create(fund); err != nil {
		return nil, fmt.Errorf("failed to create fund: %w", err)
	}
	return fund, nil
}

func (s *FundService) GetAll(userID int64) ([]*models.Fund, error) {
	funds, err := s.repo.GetAllByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch funds: %w", err)
	}
	return funds, nil
}

func (s *FundService) GetByID(id, userID int64) (*models.Fund, error) {
	fund, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fund: %w", err)
	}
	if fund == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Fund not found")
	}
	return fund, nil
}

func (s *FundService) Update(id, userID int64, req *models.UpdateFundRequest) (*models.Fund, error) {
	if req.TargetAmount < 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Target amount cannot be negative")
	}

	fund, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fund: %w", err)
	}
	if fund == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Fund not found")
	}

	fund.Name = req.Name
	fund.Description = req.Description
	fund.TargetAmount = req.TargetAmount
	if req.Currency != "" {
		fund.Currency = req.Currency
	}

	if err := s.repo.Update(fund); err != nil {
		return nil, err
	}
	return fund, nil
}

func (s *FundService) Delete(id, userID int64) error {
	if err := s.repo.Delete(id, userID); err != nil {
		if err.Error() == "fund not found" {
			return fiber.NewError(fiber.StatusNotFound, "Fund not found")
		}
		return err
	}
	return nil
}

func (s *FundService) Deposit(id, userID int64, req *models.FundTransactionRequest) (*models.Fund, error) {
	if req.Amount <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Amount must be greater than 0")
	}

	fund, err := s.repo.Deposit(id, userID, req.Amount)
	if err != nil {
		if err.Error() == "fund not found" {
			return nil, fiber.NewError(fiber.StatusNotFound, "Fund not found")
		}
		return nil, fmt.Errorf("failed to deposit: %w", err)
	}
	return fund, nil
}

func (s *FundService) Withdraw(id, userID int64, req *models.FundTransactionRequest) (*models.Fund, error) {
	if req.Amount <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Amount must be greater than 0")
	}

	fund, err := s.repo.Withdraw(id, userID, req.Amount)
	if err != nil {
		switch err.Error() {
		case "fund not found":
			return nil, fiber.NewError(fiber.StatusNotFound, "Fund not found")
		case "insufficient balance":
			return nil, fiber.NewError(fiber.StatusBadRequest, "Insufficient fund balance")
		}
		return nil, fmt.Errorf("failed to withdraw: %w", err)
	}
	return fund, nil
}
