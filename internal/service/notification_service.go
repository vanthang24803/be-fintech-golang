package service

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/push"
	"github.com/maynguyen24/sever/pkg/queue"
	"github.com/maynguyen24/sever/pkg/snowflake"
)

// NotificationRepository defines the DB contract for this service
type NotificationRepository interface {
	GetByUserID(userID int64, filter models.NotificationFilter) ([]*models.Notification, error)
	GetUnreadCount(userID int64) (int, error)
	MarkAsRead(userID int64, ids []int64) error
	Delete(userID int64, id int64) error
	Create(notif *models.Notification) error
}

type DeviceRepositoryForNotify interface {
	GetPushTokensByUserID(userID int64) ([]string, error)
}

type NotificationService struct {
	repo       NotificationRepository
	deviceRepo DeviceRepositoryForNotify
	push       push.PushClient
	queue      *queue.Client // Add queue client
}

func NewNotificationService(repo NotificationRepository, deviceRepo DeviceRepositoryForNotify, push push.PushClient, queue *queue.Client) *NotificationService {
	return &NotificationService{
		repo:       repo,
		deviceRepo: deviceRepo,
		push:       push,
		queue:      queue,
	}
}

// Create persists a notification and enqueues an asynchronous push delivery task
func (s *NotificationService) Create(notif *models.Notification) error {
	// 1. Persist to DB
	if err := s.repo.Create(notif); err != nil {
		return fmt.Errorf("failed to persist notification: %w", err)
	}

	// 2. Enqueue Push Task (Instead of goroutine)
	if s.queue != nil {
		_ = s.queue.EnqueueSendPush(queue.PushPayload{
			UserID: notif.UserID,
			Title:  notif.Title,
			Body:   notif.Body,
		})
	}

	return nil
}

// PushOnly handles local push delivery (called by worker)
func (s *NotificationService) PushOnly(ctx context.Context, userID int64, title, body string) error {
	tokens, err := s.deviceRepo.GetPushTokensByUserID(userID)
	if err != nil || len(tokens) == 0 {
		return nil
	}

	for _, token := range tokens {
		_ = s.push.SendPush(ctx, token, title, body, nil)
	}

	return nil
}

func (s *NotificationService) GetList(userID int64, filter models.NotificationFilter) ([]*models.Notification, error) {
	notifications, err := s.repo.GetByUserID(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}
	return notifications, nil
}

func (s *NotificationService) GetUnreadCount(userID int64) (int, error) {
	count, err := s.repo.GetUnreadCount(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch unread count: %w", err)
	}
	return count, nil
}

func (s *NotificationService) MarkRead(userID int64, req *models.MarkReadRequest) error {
	if len(req.IDs) == 0 {
		return nil
	}
	if err := s.repo.MarkAsRead(userID, req.IDs); err != nil {
		return fmt.Errorf("failed to mark notifications as read: %w", err)
	}
	return nil
}

func (s *NotificationService) Delete(userID int64, id int64) error {
	if err := s.repo.Delete(userID, id); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// Notifier defines the interface for creating notifications (DB + Push)
type Notifier interface {
	Create(notif *models.Notification) error
}

// SavingsGoalRepository defines the DB contract for the savings goal service
type SavingsGoalRepository interface {
	CreateGoal(goal *models.SavingsGoal) error
	GetGoalByID(id int64) (*models.SavingsGoal, error)
	ListGoals(userID int64) ([]models.SavingsGoal, error)
	UpdateGoalAmount(goalID int64, amount float64) error
	CreateContribution(c *models.GoalContribution) error
	GetContributionsByGoal(goalID int64) ([]models.GoalContribution, error)
	DeleteGoal(id int64) error
}

type SavingsGoalService struct {
	repo         SavingsGoalRepository
	fundRepo     FundRepository
	notification Notifier
}

func NewSavingsGoalService(repo SavingsGoalRepository, fundRepo FundRepository, notification Notifier, queue *queue.Client) *SavingsGoalService {
	return &SavingsGoalService{
		repo:         repo,
		fundRepo:     fundRepo,
		notification: notification,
	}
}

func (s *SavingsGoalService) Create(userID int64, req *models.CreateGoalRequest) (*models.SavingsGoal, error) {
	if req.TargetAmount <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Target amount must be greater than 0")
	}

	goal := &models.SavingsGoal{
		ID:           snowflake.GenerateID(),
		UserID:       userID,
		Name:         req.Name,
		TargetAmount: req.TargetAmount,
		TargetDate:   req.TargetDate,
		Status:       "active",
	}

	if err := s.repo.CreateGoal(goal); err != nil {
		return nil, fmt.Errorf("failed to create savings goal: %w", err)
	}
	return goal, nil
}

func (s *SavingsGoalService) List(userID int64) ([]models.SavingsGoal, error) {
	goals, err := s.repo.ListGoals(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list savings goals: %w", err)
	}

	// Calculate percentage for each
	for i := range goals {
		if goals[i].TargetAmount > 0 {
			goals[i].ProgressPercentage = (goals[i].CurrentAmount / goals[i].TargetAmount) * 100
		}
	}

	return goals, nil
}

func (s *SavingsGoalService) GetDetail(id int64, userID int64) (*models.GoalResponse, error) {
	goal, err := s.repo.GetGoalByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch goal: %w", err)
	}
	if goal == nil || goal.UserID != userID {
		return nil, fiber.NewError(fiber.StatusNotFound, "Goal not found")
	}

	if goal.TargetAmount > 0 {
		goal.ProgressPercentage = (goal.CurrentAmount / goal.TargetAmount) * 100
	}

	contributions, err := s.repo.GetContributionsByGoal(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contributions: %w", err)
	}

	return &models.GoalResponse{
		Goal:          goal,
		Contributions: contributions,
	}, nil
}

func (s *SavingsGoalService) Contribute(userID int64, req *models.GoalContributeRequest) (*models.SavingsGoal, error) {
	if req.Amount <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Contribution amount must be greater than 0")
	}

	goal, err := s.repo.GetGoalByID(req.GoalID)
	if err != nil || goal == nil || goal.UserID != userID {
		return nil, fiber.NewError(fiber.StatusNotFound, "Savings goal not found")
	}

	// 1. Withdraw from Fund
	_, err = s.fundRepo.Withdraw(req.FundID, userID, req.Amount)
	if err != nil {
		return nil, err // Pass through 'insufficient balance' or 'not found' errors
	}

	// 2. Update Goal Amount
	if err := s.repo.UpdateGoalAmount(req.GoalID, req.Amount); err != nil {
		return nil, fmt.Errorf("failed to update goal balance: %w", err)
	}

	// 3. Log Contribution
	contribution := &models.GoalContribution{
		ID:     snowflake.GenerateID(),
		GoalID: req.GoalID,
		FundID: &req.FundID,
		Amount: req.Amount,
		Type:   "deposit",
		Notes:  req.Notes,
	}
	if err := s.repo.CreateContribution(contribution); err != nil {
		return nil, fmt.Errorf("failed to log contribution: %w", err)
	}

	// 4. Reload goal for notification check
	updatedGoal, _ := s.repo.GetGoalByID(req.GoalID)

	// Check thresholds for notifications
	prevPercent := (goal.CurrentAmount / goal.TargetAmount) * 100
	currPercent := (updatedGoal.CurrentAmount / updatedGoal.TargetAmount) * 100

	s.checkGoalNotifications(updatedGoal, prevPercent, currPercent)

	return updatedGoal, nil
}

func (s *SavingsGoalService) Withdraw(userID int64, req *models.GoalWithdrawRequest) (*models.SavingsGoal, error) {
	if req.Amount <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Withdrawal amount must be greater than 0")
	}

	goal, err := s.repo.GetGoalByID(req.GoalID)
	if err != nil || goal == nil || goal.UserID != userID {
		return nil, fiber.NewError(fiber.StatusNotFound, "Savings goal not found")
	}

	if goal.CurrentAmount < req.Amount {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Insufficient goal balance")
	}

	// 1. Return to Fund
	_, err = s.fundRepo.Deposit(req.FundID, userID, req.Amount)
	if err != nil {
		return nil, err
	}

	// 2. Update Goal Amount (negative)
	if err := s.repo.UpdateGoalAmount(req.GoalID, -req.Amount); err != nil {
		return nil, fmt.Errorf("failed to update goal balance: %w", err)
	}

	// 3. Log Withdrawal
	contribution := &models.GoalContribution{
		ID:     snowflake.GenerateID(),
		GoalID: req.GoalID,
		FundID: &req.FundID,
		Amount: req.Amount,
		Type:   "withdrawal",
		Notes:  "Return to fund",
	}
	if err := s.repo.CreateContribution(contribution); err != nil {
		return nil, fmt.Errorf("failed to log withdrawal: %w", err)
	}

	updatedGoal, _ := s.repo.GetGoalByID(req.GoalID)
	return updatedGoal, nil
}

func (s *SavingsGoalService) checkGoalNotifications(goal *models.SavingsGoal, prevVal, currVal float64) {
	var nType, title, body string

	if prevVal < 100 && currVal >= 100 {
		nType = "GOAL_COMPLETED"
		title = "Goal Reached! 🎉"
		body = fmt.Sprintf("Congratulations! You have successfully reached your goal: %s.", goal.Name)
	} else if prevVal < 80 && currVal >= 80 {
		nType = "GOAL_NEAR"
		title = "Almost there! 🚀"
		body = fmt.Sprintf("You are 80%% focused on your goal: %s.", goal.Name)
	} else if prevVal < 50 && currVal >= 50 {
		nType = "GOAL_MID"
		title = "Halfway Point! 📈"
		body = fmt.Sprintf("You have reached 50%% of your goal: %s.", goal.Name)
	}

	if nType != "" {
		_ = s.notification.Create(&models.Notification{
			UserID:   goal.UserID,
			Source:   "SAVINGS_GOAL",
			SourceID: &goal.ID,
			Type:     models.NotificationType(nType),
			Title:    title,
			Body:     body,
		})
	}
}
