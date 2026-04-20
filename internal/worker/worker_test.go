package worker

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/internal/service"
	"github.com/maynguyen24/sever/pkg/push"
	"github.com/maynguyen24/sever/pkg/queue"
)

type workerDeviceRepoStub struct {
	tokens []string
	err    error
}

func (s *workerDeviceRepoStub) GetPushTokensByUserID(context.Context, int64) ([]string, error) {
	return s.tokens, s.err
}

type workerNotifyRepoStub struct{}

func (s *workerNotifyRepoStub) GetByUserID(context.Context, int64, models.NotificationFilter) ([]*models.Notification, error) {
	return nil, nil
}
func (s *workerNotifyRepoStub) GetUnreadCount(context.Context, int64) (int, error) { return 0, nil }
func (s *workerNotifyRepoStub) MarkAsRead(context.Context, int64, []int64) error    { return nil }
func (s *workerNotifyRepoStub) Delete(context.Context, int64, int64) error           { return nil }
func (s *workerNotifyRepoStub) Create(context.Context, *models.Notification) error   { return nil }

type workerBudgetRepoStub struct {
	getByUserIDFn        func(context.Context, int64) ([]*models.Budget, error)
	calculateSpendingFn  func(context.Context, int64, *int64, time.Time, time.Time) (float64, error)
}

func (s *workerBudgetRepoStub) GetByUserID(ctx context.Context, userID int64) ([]*models.Budget, error) {
	if s.getByUserIDFn != nil {
		return s.getByUserIDFn(ctx, userID)
	}
	return nil, nil
}

func (s *workerBudgetRepoStub) CalculateSpending(ctx context.Context, userID int64, categoryID *int64, start, end time.Time) (float64, error) {
	if s.calculateSpendingFn != nil {
		return s.calculateSpendingFn(ctx, userID, categoryID, start, end)
	}
	return 0, nil
}

type workerTxRepoStub struct{}

func (s *workerTxRepoStub) Create(context.Context, *models.Transaction) error { return nil }
func (s *workerTxRepoStub) GetAllByUserID(context.Context, int64, models.TransactionFilter) ([]*models.TransactionDetail, error) {
	return nil, nil
}
func (s *workerTxRepoStub) GetByID(context.Context, int64, int64) (*models.TransactionDetail, error) {
	return nil, nil
}
func (s *workerTxRepoStub) GetRawByID(context.Context, int64, int64) (*models.Transaction, error) {
	return nil, nil
}
func (s *workerTxRepoStub) Update(context.Context, *models.Transaction, *models.Transaction) error { return nil }
func (s *workerTxRepoStub) Delete(context.Context, int64, int64) error                             { return nil }

type workerNotificationRecorder struct {
	created []*models.Notification
	err     error
}

func (s *workerNotificationRecorder) Create(_ context.Context, notif *models.Notification) error {
	if s.err != nil {
		return s.err
	}
	s.created = append(s.created, notif)
	return nil
}

type workerPushStub struct {
	sent []string
}

func (s *workerPushStub) SendPush(_ context.Context, token, title, body string, data map[string]string) error {
	s.sent = append(s.sent, token+"|"+title+"|"+body)
	return nil
}

func (s *workerPushStub) SendToTopic(context.Context, string, string, string, map[string]string) error {
	return nil
}

func TestHandlers_HandleSendPushTask(t *testing.T) {
	t.Parallel()

	pushClient := &workerPushStub{}
	notifySvc := service.NewNotificationService(&workerNotifyRepoStub{}, &workerDeviceRepoStub{tokens: []string{"t1", "t2"}}, pushClient, nil)
	handlers := NewHandlers(notifySvc, nil)

	payload, err := json.Marshal(queue.PushPayload{UserID: 42, Title: "Hi", Body: "Body"})
	if err != nil {
		t.Fatalf("Marshal payload: %v", err)
	}

	if err := handlers.HandleSendPushTask(context.Background(), asynq.NewTask(queue.TypeSendPush, payload)); err != nil {
		t.Fatalf("HandleSendPushTask() error = %v", err)
	}
	if len(pushClient.sent) != 2 {
		t.Fatalf("expected push to be sent to 2 tokens, got %d", len(pushClient.sent))
	}
}

func TestHandlers_HandleSendPushTask_InvalidPayload(t *testing.T) {
	t.Parallel()

	handlers := NewHandlers(service.NewNotificationService(&workerNotifyRepoStub{}, &workerDeviceRepoStub{}, &push.MockPushClient{}, nil), nil)
	err := handlers.HandleSendPushTask(context.Background(), asynq.NewTask(queue.TypeSendPush, []byte("{")))
	if !errors.Is(err, asynq.SkipRetry) {
		t.Fatalf("expected skip retry error, got %v", err)
	}
}

func TestHandlers_HandleCheckBudgetTask(t *testing.T) {
	t.Parallel()

	categoryID := int64(9)
	recorder := &workerNotificationRecorder{}
	txSvc := service.NewTransactionService(
		&workerTxRepoStub{},
		&workerBudgetRepoStub{
			getByUserIDFn: func(context.Context, int64) ([]*models.Budget, error) {
				return []*models.Budget{{
					ID:         1,
					UserID:     42,
					CategoryID: &categoryID,
					Amount:     100,
					IsActive:   true,
				}}, nil
			},
			calculateSpendingFn: func(context.Context, int64, *int64, time.Time, time.Time) (float64, error) {
				return 90, nil
			},
		},
		recorder,
		nil,
	)
	handlers := NewHandlers(nil, txSvc)

	payload, err := json.Marshal(queue.BudgetCheckPayload{UserID: 42, CategoryID: &categoryID})
	if err != nil {
		t.Fatalf("Marshal payload: %v", err)
	}

	if err := handlers.HandleCheckBudgetTask(context.Background(), asynq.NewTask(queue.TypeCheckBudget, payload)); err != nil {
		t.Fatalf("HandleCheckBudgetTask() error = %v", err)
	}
	if len(recorder.created) != 1 {
		t.Fatalf("expected 1 notification created, got %d", len(recorder.created))
	}
}

func TestHandlers_HandleCheckBudgetTask_Errors(t *testing.T) {
	t.Parallel()

	handlers := NewHandlers(nil, service.NewTransactionService(&workerTxRepoStub{}, &workerBudgetRepoStub{
		getByUserIDFn: func(context.Context, int64) ([]*models.Budget, error) {
			return nil, errors.New("budget lookup failed")
		},
	}, &workerNotificationRecorder{}, nil))

	err := handlers.HandleCheckBudgetTask(context.Background(), asynq.NewTask(queue.TypeCheckBudget, []byte("{")))
	if !errors.Is(err, asynq.SkipRetry) {
		t.Fatalf("expected skip retry for invalid payload, got %v", err)
	}

	payload, _ := json.Marshal(queue.BudgetCheckPayload{UserID: 42})
	err = handlers.HandleCheckBudgetTask(context.Background(), asynq.NewTask(queue.TypeCheckBudget, payload))
	if err == nil || err.Error() != "failed to check budgets: budget lookup failed" {
		t.Fatalf("expected wrapped budget error, got %v", err)
	}
}

func TestNewWorkerRegisterHandlerStartStop(t *testing.T) {
	t.Parallel()

	worker := NewWorker("127.0.0.1:0", "", 2)
	if worker == nil || worker.server == nil || worker.mux == nil {
		t.Fatal("expected worker server and mux to be initialized")
	}

	worker.RegisterHandler(queue.TypeSendPush, func(ctx context.Context, t *asynq.Task) error {
		return nil
	})

	if err := worker.Start(); err != nil {
		t.Fatalf("expected worker start to return nil, got %v", err)
	}
	worker.Stop()
}
