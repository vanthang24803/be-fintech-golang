package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/maynguyen24/sever/internal/service"
	"github.com/maynguyen24/sever/pkg/queue"
)

type Handlers struct {
	notifyService *service.NotificationService
	txService     *service.TransactionService
}

func NewHandlers(notifyService *service.NotificationService, txService *service.TransactionService) *Handlers {
	return &Handlers{
		notifyService: notifyService,
		txService:     txService,
	}
}

func (h *Handlers) HandleSendPushTask(ctx context.Context, t *asynq.Task) error {
	var p queue.PushPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Worker: Sending push notification to UserID %d", p.UserID)
	return h.notifyService.PushOnly(ctx, p.UserID, p.Title, p.Body)
}

func (h *Handlers) HandleCheckBudgetTask(ctx context.Context, t *asynq.Task) error {
	var p queue.BudgetCheckPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Worker: Checking budget for UserID %d", p.UserID)
	if err := h.txService.CheckBudgets(ctx, p.UserID, p.CategoryID); err != nil {
		return fmt.Errorf("failed to check budgets: %w", err)
	}
	return nil
}
