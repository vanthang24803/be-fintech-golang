package service

import (
	"context"
	"fmt"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
	"github.com/maynguyen24/sever/pkg/push"
	"github.com/maynguyen24/sever/pkg/queue"
	"github.com/maynguyen24/sever/pkg/snowflake"
)

// NotificationRepository defines the DB contract for this service
type NotificationRepository interface {
	GetByUserID(ctx context.Context, userID int64, filter models.NotificationFilter) ([]*models.Notification, error)
	GetUnreadCount(ctx context.Context, userID int64) (int, error)
	MarkAsRead(ctx context.Context, userID int64, ids []int64) error
	Delete(ctx context.Context, userID int64, id int64) error
	Create(ctx context.Context, notif *models.Notification) error
}

type DeviceRepositoryForNotify interface {
	GetPushTokensByUserID(ctx context.Context, userID int64) ([]string, error)
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
func (s *NotificationService) Create(ctx context.Context, notif *models.Notification) error {
	// 1. Persist to DB
	if err := s.repo.Create(ctx, notif); err != nil {
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
	tokens, err := s.deviceRepo.GetPushTokensByUserID(ctx, userID)
	if err != nil || len(tokens) == 0 {
		return nil
	}

	for _, token := range tokens {
		_ = s.push.SendPush(ctx, token, title, body, nil)
	}

	return nil
}

func (s *NotificationService) GetList(ctx context.Context, userID int64, filter models.NotificationFilter) ([]*models.Notification, error) {
	notifications, err := s.repo.GetByUserID(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}
	return notifications, nil
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID int64) (int, error) {
	count, err := s.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch unread count: %w", err)
	}
	return count, nil
}

func (s *NotificationService) MarkRead(ctx context.Context, userID int64, req *models.MarkReadRequest) error {
	if len(req.IDs) == 0 {
		return nil
	}
	if err := s.repo.MarkAsRead(ctx, userID, req.IDs); err != nil {
		return fmt.Errorf("failed to mark notifications as read: %w", err)
	}
	return nil
}

func (s *NotificationService) Delete(ctx context.Context, userID int64, id int64) error {
	if err := s.repo.Delete(ctx, userID, id); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// Notifier defines the interface for creating notifications (DB + Push)
type Notifier interface {
	Create(ctx context.Context, notif *models.Notification) error
}

// SavingsGoalRepository defines the DB contract for the savings goal service
type SavingsGoalRepository interface {
	CreateGoal(ctx context.Context, goal *models.SavingsGoal) error
	GetGoalByID(ctx context.Context, id int64) (*models.SavingsGoal, error)
	ListGoals(ctx context.Context, userID int64) ([]models.SavingsGoal, error)
	UpdateGoalAmount(ctx context.Context, goalID int64, amount float64) error
	CreateContribution(ctx context.Context, c *models.GoalContribution) error
	GetContributionsByGoal(ctx context.Context, goalID int64) ([]models.GoalContribution, error)
	DeleteGoal(ctx context.Context, id int64) error
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

func (s *SavingsGoalService) Create(ctx context.Context, userID int64, req *models.CreateGoalRequest) (*models.SavingsGoal, error) {
	if req.TargetAmount <= 0 {
		return nil, fmt.Errorf("%w: Target amount must be greater than 0", apperr.ErrInvalidInput)
	}

	goal := &models.SavingsGoal{
		ID:           snowflake.GenerateID(),
		UserID:       userID,
		Name:         req.Name,
		TargetAmount: req.TargetAmount,
		TargetDate:   req.TargetDate,
		Status:       "active",
	}

	if err := s.repo.CreateGoal(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to create savings goal: %w", err)
	}
	return goal, nil
}

func (s *SavingsGoalService) List(ctx context.Context, userID int64) ([]models.SavingsGoal, error) {
	goals, err := s.repo.ListGoals(ctx, userID)
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

func (s *SavingsGoalService) GetDetail(ctx context.Context, id int64, userID int64) (*models.GoalResponse, error) {
	goal, err := s.repo.GetGoalByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch goal: %w", err)
	}
	if goal == nil || goal.UserID != userID {
		return nil, fmt.Errorf("%w: Goal not found", apperr.ErrNotFound)
	}

	if goal.TargetAmount > 0 {
		goal.ProgressPercentage = (goal.CurrentAmount / goal.TargetAmount) * 100
	}

	contributions, err := s.repo.GetContributionsByGoal(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contributions: %w", err)
	}

	return &models.GoalResponse{
		Goal:          goal,
		Contributions: contributions,
	}, nil
}

func (s *SavingsGoalService) Contribute(ctx context.Context, userID int64, req *models.GoalContributeRequest) (*models.SavingsGoal, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("%w: Contribution amount must be greater than 0", apperr.ErrInvalidInput)
	}

	goal, err := s.repo.GetGoalByID(ctx, req.GoalID)
	if err != nil || goal == nil || goal.UserID != userID {
		return nil, fmt.Errorf("%w: Savings goal not found", apperr.ErrNotFound)
	}

	// 1. Withdraw from Fund
	_, err = s.fundRepo.Withdraw(ctx, req.FundID, userID, req.Amount)
	if err != nil {
		return nil, err // Pass through 'insufficient balance' or 'not found' errors
	}

	// 2. Update Goal Amount
	if err := s.repo.UpdateGoalAmount(ctx, req.GoalID, req.Amount); err != nil {
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
	if err := s.repo.CreateContribution(ctx, contribution); err != nil {
		return nil, fmt.Errorf("failed to log contribution: %w", err)
	}

	// 4. Reload goal for notification check
	updatedGoal, _ := s.repo.GetGoalByID(ctx, req.GoalID)

	// Check thresholds for notifications
	prevPercent := (goal.CurrentAmount / goal.TargetAmount) * 100
	currPercent := (updatedGoal.CurrentAmount / updatedGoal.TargetAmount) * 100

	s.checkGoalNotifications(ctx, updatedGoal, prevPercent, currPercent)

	return updatedGoal, nil
}

func (s *SavingsGoalService) Withdraw(ctx context.Context, userID int64, req *models.GoalWithdrawRequest) (*models.SavingsGoal, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("%w: Withdrawal amount must be greater than 0", apperr.ErrInvalidInput)
	}

	goal, err := s.repo.GetGoalByID(ctx, req.GoalID)
	if err != nil || goal == nil || goal.UserID != userID {
		return nil, fmt.Errorf("%w: Savings goal not found", apperr.ErrNotFound)
	}

	if goal.CurrentAmount < req.Amount {
		return nil, fmt.Errorf("%w: Insufficient goal balance", apperr.ErrInvalidInput)
	}

	// 1. Return to Fund
	_, err = s.fundRepo.Deposit(ctx, req.FundID, userID, req.Amount)
	if err != nil {
		return nil, err
	}

	// 2. Update Goal Amount (negative)
	if err := s.repo.UpdateGoalAmount(ctx, req.GoalID, -req.Amount); err != nil {
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
	if err := s.repo.CreateContribution(ctx, contribution); err != nil {
		return nil, fmt.Errorf("failed to log withdrawal: %w", err)
	}

	updatedGoal, _ := s.repo.GetGoalByID(ctx, req.GoalID)
	return updatedGoal, nil
}

func (s *SavingsGoalService) checkGoalNotifications(ctx context.Context, goal *models.SavingsGoal, prevVal, currVal float64) {
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
		_ = s.notification.Create(ctx, &models.Notification{
			UserID:   goal.UserID,
			Source:   "SAVINGS_GOAL",
			SourceID: &goal.ID,
			Type:     models.NotificationType(nType),
			Title:    title,
			Body:     body,
		})
	}
}
