package service

import (
	"context"
	"errors"
	"testing"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
	"github.com/maynguyen24/sever/pkg/push"
)

// --- stubs for NotificationService ---

type stubFullNotifRepo struct {
	getByUserIDFn    func(context.Context, int64, models.NotificationFilter) ([]*models.Notification, error)
	getUnreadCountFn func(context.Context, int64) (int, error)
	markAsReadFn     func(context.Context, int64, []int64) error
	deleteNotifFn    func(context.Context, int64, int64) error
	createNotifFn    func(context.Context, *models.Notification) error
}

func (s *stubFullNotifRepo) GetByUserID(ctx context.Context, userID int64, filter models.NotificationFilter) ([]*models.Notification, error) {
	if s.getByUserIDFn != nil {
		return s.getByUserIDFn(ctx, userID, filter)
	}
	return nil, nil
}

func (s *stubFullNotifRepo) GetUnreadCount(ctx context.Context, userID int64) (int, error) {
	if s.getUnreadCountFn != nil {
		return s.getUnreadCountFn(ctx, userID)
	}
	return 0, nil
}

func (s *stubFullNotifRepo) MarkAsRead(ctx context.Context, userID int64, ids []int64) error {
	if s.markAsReadFn != nil {
		return s.markAsReadFn(ctx, userID, ids)
	}
	return nil
}

func (s *stubFullNotifRepo) Delete(ctx context.Context, userID int64, id int64) error {
	if s.deleteNotifFn != nil {
		return s.deleteNotifFn(ctx, userID, id)
	}
	return nil
}

func (s *stubFullNotifRepo) Create(ctx context.Context, notif *models.Notification) error {
	if s.createNotifFn != nil {
		return s.createNotifFn(ctx, notif)
	}
	return nil
}

type stubNotifDeviceRepo struct {
	getTokensFn func(context.Context, int64) ([]string, error)
}

func (s *stubNotifDeviceRepo) GetPushTokensByUserID(ctx context.Context, userID int64) ([]string, error) {
	if s.getTokensFn != nil {
		return s.getTokensFn(ctx, userID)
	}
	return nil, nil
}

type stubPushClient struct {
	sendFn func(context.Context, string, string, string, map[string]string) error
}

func (s *stubPushClient) SendPush(ctx context.Context, token, title, body string, data map[string]string) error {
	if s.sendFn != nil {
		return s.sendFn(ctx, token, title, body, data)
	}
	return nil
}

func (s *stubPushClient) SendToTopic(ctx context.Context, topic, title, body string, data map[string]string) error {
	return nil
}

// --- NotificationService tests ---

func TestNotificationService_Create(t *testing.T) {
	t.Parallel()

	var persisted *models.Notification
	repo := &stubFullNotifRepo{
		createNotifFn: func(ctx context.Context, notif *models.Notification) error {
			persisted = notif
			return nil
		},
	}

	svc := NewNotificationService(repo, &stubNotifDeviceRepo{}, &push.MockPushClient{}, nil)
	err := svc.Create(context.Background(), &models.Notification{
		UserID: 1,
		Title:  "Test",
		Body:   "Test body",
		Type:   models.NotifInfo,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if persisted == nil {
		t.Fatal("expected notification to be persisted")
	}
}

func TestNotificationService_Create_RepoError(t *testing.T) {
	t.Parallel()

	repo := &stubFullNotifRepo{
		createNotifFn: func(ctx context.Context, notif *models.Notification) error {
			return errors.New("db error")
		},
	}

	svc := NewNotificationService(repo, &stubNotifDeviceRepo{}, &push.MockPushClient{}, nil)
	err := svc.Create(context.Background(), &models.Notification{})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestNotificationService_PushOnly_NoTokens(t *testing.T) {
	t.Parallel()

	svc := NewNotificationService(&stubFullNotifRepo{}, &stubNotifDeviceRepo{
		getTokensFn: func(ctx context.Context, userID int64) ([]string, error) {
			return nil, nil
		},
	}, &push.MockPushClient{}, nil)

	if err := svc.PushOnly(context.Background(), 1, "Title", "Body"); err != nil {
		t.Fatalf("PushOnly returned error: %v", err)
	}
}

func TestNotificationService_PushOnly_WithTokens(t *testing.T) {
	t.Parallel()

	pushed := 0
	svc := NewNotificationService(&stubFullNotifRepo{}, &stubNotifDeviceRepo{
		getTokensFn: func(ctx context.Context, userID int64) ([]string, error) {
			return []string{"token1", "token2"}, nil
		},
	}, &stubPushClient{
		sendFn: func(ctx context.Context, token, title, body string, data map[string]string) error {
			pushed++
			return nil
		},
	}, nil)

	if err := svc.PushOnly(context.Background(), 1, "Title", "Body"); err != nil {
		t.Fatalf("PushOnly returned error: %v", err)
	}
	if pushed != 2 {
		t.Fatalf("expected 2 push calls, got %d", pushed)
	}
}

func TestNotificationService_GetList(t *testing.T) {
	t.Parallel()

	notifs := []*models.Notification{{ID: 1, UserID: 42}, {ID: 2, UserID: 42}}
	svc := NewNotificationService(&stubFullNotifRepo{
		getByUserIDFn: func(ctx context.Context, userID int64, filter models.NotificationFilter) ([]*models.Notification, error) {
			return notifs, nil
		},
	}, &stubNotifDeviceRepo{}, &push.MockPushClient{}, nil)

	got, err := svc.GetList(context.Background(), 42, models.NotificationFilter{})
	if err != nil {
		t.Fatalf("GetList returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 notifications, got %d", len(got))
	}
}

func TestNotificationService_GetUnreadCount(t *testing.T) {
	t.Parallel()

	svc := NewNotificationService(&stubFullNotifRepo{
		getUnreadCountFn: func(ctx context.Context, userID int64) (int, error) {
			return 5, nil
		},
	}, &stubNotifDeviceRepo{}, &push.MockPushClient{}, nil)

	count, err := svc.GetUnreadCount(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetUnreadCount returned error: %v", err)
	}
	if count != 5 {
		t.Fatalf("expected count 5, got %d", count)
	}
}

func TestNotificationService_MarkRead_EmptyIDs(t *testing.T) {
	t.Parallel()

	svc := NewNotificationService(&stubFullNotifRepo{}, &stubNotifDeviceRepo{}, &push.MockPushClient{}, nil)
	if err := svc.MarkRead(context.Background(), 42, &models.MarkReadRequest{IDs: nil}); err != nil {
		t.Fatalf("MarkRead with empty IDs returned error: %v", err)
	}
}

func TestNotificationService_MarkRead(t *testing.T) {
	t.Parallel()

	var markedIDs []int64
	svc := NewNotificationService(&stubFullNotifRepo{
		markAsReadFn: func(ctx context.Context, userID int64, ids []int64) error {
			markedIDs = ids
			return nil
		},
	}, &stubNotifDeviceRepo{}, &push.MockPushClient{}, nil)

	if err := svc.MarkRead(context.Background(), 42, &models.MarkReadRequest{IDs: []int64{1, 2, 3}}); err != nil {
		t.Fatalf("MarkRead returned error: %v", err)
	}
	if len(markedIDs) != 3 {
		t.Fatalf("expected 3 IDs marked, got %d", len(markedIDs))
	}
}

func TestNotificationService_Delete(t *testing.T) {
	t.Parallel()

	deleted := false
	svc := NewNotificationService(&stubFullNotifRepo{
		deleteNotifFn: func(ctx context.Context, userID int64, id int64) error {
			deleted = true
			return nil
		},
	}, &stubNotifDeviceRepo{}, &push.MockPushClient{}, nil)

	if err := svc.Delete(context.Background(), 42, 1); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if !deleted {
		t.Fatal("expected delete to be called")
	}
}

// --- stubs for SavingsGoalService ---

type stubSavingsGoalRepo struct {
	createGoalFn             func(context.Context, *models.SavingsGoal) error
	getGoalByIDFn            func(context.Context, int64) (*models.SavingsGoal, error)
	listGoalsFn              func(context.Context, int64) ([]models.SavingsGoal, error)
	updateGoalAmountFn       func(context.Context, int64, float64) error
	createContributionFn     func(context.Context, *models.GoalContribution) error
	getContributionsByGoalFn func(context.Context, int64) ([]models.GoalContribution, error)
	deleteGoalFn             func(context.Context, int64) error
}

func (s *stubSavingsGoalRepo) CreateGoal(ctx context.Context, goal *models.SavingsGoal) error {
	if s.createGoalFn != nil {
		return s.createGoalFn(ctx, goal)
	}
	return nil
}

func (s *stubSavingsGoalRepo) GetGoalByID(ctx context.Context, id int64) (*models.SavingsGoal, error) {
	if s.getGoalByIDFn != nil {
		return s.getGoalByIDFn(ctx, id)
	}
	return nil, nil
}

func (s *stubSavingsGoalRepo) ListGoals(ctx context.Context, userID int64) ([]models.SavingsGoal, error) {
	if s.listGoalsFn != nil {
		return s.listGoalsFn(ctx, userID)
	}
	return nil, nil
}

func (s *stubSavingsGoalRepo) UpdateGoalAmount(ctx context.Context, goalID int64, amount float64) error {
	if s.updateGoalAmountFn != nil {
		return s.updateGoalAmountFn(ctx, goalID, amount)
	}
	return nil
}

func (s *stubSavingsGoalRepo) CreateContribution(ctx context.Context, c *models.GoalContribution) error {
	if s.createContributionFn != nil {
		return s.createContributionFn(ctx, c)
	}
	return nil
}

func (s *stubSavingsGoalRepo) GetContributionsByGoal(ctx context.Context, goalID int64) ([]models.GoalContribution, error) {
	if s.getContributionsByGoalFn != nil {
		return s.getContributionsByGoalFn(ctx, goalID)
	}
	return nil, nil
}

func (s *stubSavingsGoalRepo) DeleteGoal(ctx context.Context, id int64) error {
	if s.deleteGoalFn != nil {
		return s.deleteGoalFn(ctx, id)
	}
	return nil
}

type stubSavingsFundRepo struct {
	depositFn  func(context.Context, int64, int64, float64) (*models.Fund, error)
	withdrawFn func(context.Context, int64, int64, float64) (*models.Fund, error)
}

func (s *stubSavingsFundRepo) Create(ctx context.Context, fund *models.Fund) error { return nil }
func (s *stubSavingsFundRepo) GetAllByUserID(ctx context.Context, userID int64) ([]*models.Fund, error) {
	return nil, nil
}
func (s *stubSavingsFundRepo) GetByID(ctx context.Context, id, userID int64) (*models.Fund, error) {
	return nil, nil
}
func (s *stubSavingsFundRepo) Update(ctx context.Context, fund *models.Fund) error { return nil }
func (s *stubSavingsFundRepo) Delete(ctx context.Context, id, userID int64) error  { return nil }
func (s *stubSavingsFundRepo) Deposit(ctx context.Context, id, userID int64, amount float64) (*models.Fund, error) {
	if s.depositFn != nil {
		return s.depositFn(ctx, id, userID, amount)
	}
	return &models.Fund{ID: id}, nil
}
func (s *stubSavingsFundRepo) Withdraw(ctx context.Context, id, userID int64, amount float64) (*models.Fund, error) {
	if s.withdrawFn != nil {
		return s.withdrawFn(ctx, id, userID, amount)
	}
	return &models.Fund{ID: id}, nil
}

type stubNotifier struct {
	createFn func(context.Context, *models.Notification) error
}

func (s *stubNotifier) Create(ctx context.Context, notif *models.Notification) error {
	if s.createFn != nil {
		return s.createFn(ctx, notif)
	}
	return nil
}

// --- SavingsGoalService tests ---

func TestSavingsGoalService_Create(t *testing.T) {
	t.Parallel()

	var created *models.SavingsGoal
	repo := &stubSavingsGoalRepo{
		createGoalFn: func(ctx context.Context, g *models.SavingsGoal) error {
			created = g
			return nil
		},
	}

	svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	got, err := svc.Create(context.Background(), 42, &models.CreateGoalRequest{
		Name:         "New Car",
		TargetAmount: 5000,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got == nil || created == nil {
		t.Fatal("expected goal to be created")
	}
	if created.Name != "New Car" || created.TargetAmount != 5000 || created.UserID != 42 {
		t.Fatalf("unexpected goal created: %+v", created)
	}
	if created.Status != "active" {
		t.Fatalf("expected status 'active', got %s", created.Status)
	}
}

func TestSavingsGoalService_Create_InvalidAmount(t *testing.T) {
	t.Parallel()

	svc := NewSavingsGoalService(&stubSavingsGoalRepo{}, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	_, err := svc.Create(context.Background(), 42, &models.CreateGoalRequest{
		Name:         "Test",
		TargetAmount: 0,
	})
	if !errors.Is(err, apperr.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestSavingsGoalService_Create_NegativeAmount(t *testing.T) {
	t.Parallel()

	svc := NewSavingsGoalService(&stubSavingsGoalRepo{}, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	_, err := svc.Create(context.Background(), 42, &models.CreateGoalRequest{
		Name:         "Test",
		TargetAmount: -100,
	})
	if !errors.Is(err, apperr.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestSavingsGoalService_List(t *testing.T) {
	t.Parallel()

	goals := []models.SavingsGoal{
		{ID: 1, UserID: 42, Name: "Car", TargetAmount: 1000, CurrentAmount: 500},
		{ID: 2, UserID: 42, Name: "Vacation", TargetAmount: 0, CurrentAmount: 100},
	}
	repo := &stubSavingsGoalRepo{
		listGoalsFn: func(ctx context.Context, userID int64) ([]models.SavingsGoal, error) {
			return goals, nil
		},
	}

	svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	got, err := svc.List(context.Background(), 42)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 goals, got %d", len(got))
	}
	if got[0].ProgressPercentage != 50 {
		t.Fatalf("expected progress 50, got %v", got[0].ProgressPercentage)
	}
	if got[1].ProgressPercentage != 0 {
		t.Fatalf("expected progress 0 for zero target, got %v", got[1].ProgressPercentage)
	}
}

func TestSavingsGoalService_GetDetail(t *testing.T) {
	t.Parallel()

	goal := &models.SavingsGoal{ID: 1, UserID: 42, TargetAmount: 1000, CurrentAmount: 300}
	contributions := []models.GoalContribution{{ID: 1, GoalID: 1, Amount: 300}}
	repo := &stubSavingsGoalRepo{
		getGoalByIDFn: func(ctx context.Context, id int64) (*models.SavingsGoal, error) {
			return goal, nil
		},
		getContributionsByGoalFn: func(ctx context.Context, goalID int64) ([]models.GoalContribution, error) {
			return contributions, nil
		},
	}

	svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	resp, err := svc.GetDetail(context.Background(), 1, 42)
	if err != nil {
		t.Fatalf("GetDetail returned error: %v", err)
	}
	if resp.Goal == nil {
		t.Fatal("expected goal to be returned")
	}
	if resp.Goal.ProgressPercentage != 30 {
		t.Fatalf("expected progress 30, got %v", resp.Goal.ProgressPercentage)
	}
	if len(resp.Contributions) != 1 {
		t.Fatalf("expected 1 contribution, got %d", len(resp.Contributions))
	}
}

func TestSavingsGoalService_GetDetail_NotFound(t *testing.T) {
	t.Parallel()

	repo := &stubSavingsGoalRepo{
		getGoalByIDFn: func(ctx context.Context, id int64) (*models.SavingsGoal, error) {
			return nil, nil
		},
	}

	svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	_, err := svc.GetDetail(context.Background(), 99, 42)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestSavingsGoalService_GetDetail_WrongUser(t *testing.T) {
	t.Parallel()

	repo := &stubSavingsGoalRepo{
		getGoalByIDFn: func(ctx context.Context, id int64) (*models.SavingsGoal, error) {
			return &models.SavingsGoal{ID: 1, UserID: 99}, nil
		},
	}

	svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	_, err := svc.GetDetail(context.Background(), 1, 42)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found for wrong user, got %v", err)
	}
}

func TestSavingsGoalService_Contribute(t *testing.T) {
	t.Parallel()

	goal := &models.SavingsGoal{ID: 1, UserID: 42, TargetAmount: 1000, CurrentAmount: 400}
	updatedGoal := &models.SavingsGoal{ID: 1, UserID: 42, TargetAmount: 1000, CurrentAmount: 600}

	callCount := 0
	repo := &stubSavingsGoalRepo{
		getGoalByIDFn: func(ctx context.Context, id int64) (*models.SavingsGoal, error) {
			callCount++
			if callCount == 1 {
				return goal, nil
			}
			return updatedGoal, nil
		},
		updateGoalAmountFn: func(ctx context.Context, goalID int64, amount float64) error {
			return nil
		},
		createContributionFn: func(ctx context.Context, c *models.GoalContribution) error {
			return nil
		},
	}

	svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	got, err := svc.Contribute(context.Background(), 42, &models.GoalContributeRequest{
		GoalID: 1,
		FundID: 2,
		Amount: 200,
	})
	if err != nil {
		t.Fatalf("Contribute returned error: %v", err)
	}
	if got == nil {
		t.Fatal("expected updated goal to be returned")
	}
}

func TestSavingsGoalService_Contribute_InvalidAmount(t *testing.T) {
	t.Parallel()

	svc := NewSavingsGoalService(&stubSavingsGoalRepo{}, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	_, err := svc.Contribute(context.Background(), 42, &models.GoalContributeRequest{
		GoalID: 1,
		FundID: 2,
		Amount: 0,
	})
	if !errors.Is(err, apperr.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestSavingsGoalService_Contribute_GoalNotFound(t *testing.T) {
	t.Parallel()

	repo := &stubSavingsGoalRepo{
		getGoalByIDFn: func(ctx context.Context, id int64) (*models.SavingsGoal, error) {
			return nil, nil
		},
	}

	svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	_, err := svc.Contribute(context.Background(), 42, &models.GoalContributeRequest{
		GoalID: 99,
		FundID: 2,
		Amount: 100,
	})
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestSavingsGoalService_Withdraw(t *testing.T) {
	t.Parallel()

	goal := &models.SavingsGoal{ID: 1, UserID: 42, CurrentAmount: 500}
	updatedGoal := &models.SavingsGoal{ID: 1, UserID: 42, CurrentAmount: 300}

	callCount := 0
	repo := &stubSavingsGoalRepo{
		getGoalByIDFn: func(ctx context.Context, id int64) (*models.SavingsGoal, error) {
			callCount++
			if callCount == 1 {
				return goal, nil
			}
			return updatedGoal, nil
		},
		updateGoalAmountFn: func(ctx context.Context, goalID int64, amount float64) error {
			return nil
		},
		createContributionFn: func(ctx context.Context, c *models.GoalContribution) error {
			return nil
		},
	}

	svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	got, err := svc.Withdraw(context.Background(), 42, &models.GoalWithdrawRequest{
		GoalID: 1,
		FundID: 2,
		Amount: 200,
	})
	if err != nil {
		t.Fatalf("Withdraw returned error: %v", err)
	}
	if got == nil {
		t.Fatal("expected updated goal to be returned")
	}
}

func TestSavingsGoalService_Withdraw_InvalidAmount(t *testing.T) {
	t.Parallel()

	svc := NewSavingsGoalService(&stubSavingsGoalRepo{}, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	_, err := svc.Withdraw(context.Background(), 42, &models.GoalWithdrawRequest{
		GoalID: 1,
		FundID: 2,
		Amount: 0,
	})
	if !errors.Is(err, apperr.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestSavingsGoalService_Withdraw_InsufficientBalance(t *testing.T) {
	t.Parallel()

	repo := &stubSavingsGoalRepo{
		getGoalByIDFn: func(ctx context.Context, id int64) (*models.SavingsGoal, error) {
			return &models.SavingsGoal{ID: 1, UserID: 42, CurrentAmount: 100}, nil
		},
	}

	svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	_, err := svc.Withdraw(context.Background(), 42, &models.GoalWithdrawRequest{
		GoalID: 1,
		FundID: 2,
		Amount: 200,
	})
	if !errors.Is(err, apperr.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestSavingsGoalService_Contribute_Notifications(t *testing.T) {
	t.Parallel()

	type notifCase struct {
		name      string
		prevPct   float64
		currPct   float64
		wantNotif bool
		wantType  string
	}

	cases := []notifCase{
		{"no threshold", 30, 45, false, ""},
		{"crosses 50%", 40, 55, true, "GOAL_MID"},
		{"crosses 80%", 70, 85, true, "GOAL_NEAR"},
		{"crosses 100%", 90, 105, true, "GOAL_COMPLETED"},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			targetAmount := float64(1000)
			prevAmount := tt.prevPct / 100 * targetAmount
			currAmount := tt.currPct / 100 * targetAmount
			delta := currAmount - prevAmount

			goal := &models.SavingsGoal{ID: 1, UserID: 42, TargetAmount: targetAmount, CurrentAmount: prevAmount}
			updatedGoal := &models.SavingsGoal{ID: 1, UserID: 42, TargetAmount: targetAmount, CurrentAmount: currAmount}

			callCount := 0
			repo := &stubSavingsGoalRepo{
				getGoalByIDFn: func(ctx context.Context, id int64) (*models.SavingsGoal, error) {
					callCount++
					if callCount == 1 {
						return goal, nil
					}
					return updatedGoal, nil
				},
				updateGoalAmountFn: func(ctx context.Context, goalID int64, amount float64) error {
					return nil
				},
				createContributionFn: func(ctx context.Context, c *models.GoalContribution) error {
					return nil
				},
			}

			var notifType string
			notifier := &stubNotifier{
				createFn: func(ctx context.Context, notif *models.Notification) error {
					notifType = string(notif.Type)
					return nil
				},
			}

			svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, notifier, nil)
			_, err := svc.Contribute(context.Background(), 42, &models.GoalContributeRequest{
				GoalID: 1,
				FundID: 2,
				Amount: delta,
			})
			if err != nil {
				t.Fatalf("Contribute returned error: %v", err)
			}

			if tt.wantNotif && notifType == "" {
				t.Fatalf("expected notification of type %s, got none", tt.wantType)
			}
			if !tt.wantNotif && notifType != "" {
				t.Fatalf("expected no notification, got type %s", notifType)
			}
			if tt.wantNotif && notifType != tt.wantType {
				t.Fatalf("expected notification type %s, got %s", tt.wantType, notifType)
			}
		})
	}
}
