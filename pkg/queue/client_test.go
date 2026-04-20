package queue

import "testing"

func TestNewClientAndClose(t *testing.T) {
	t.Parallel()

	client := NewClient("127.0.0.1:0", "")
	if client == nil || client.client == nil {
		t.Fatal("expected asynq client to be initialized")
	}
	if err := client.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestEnqueueMethodsReturnConnectionErrorWithoutRedis(t *testing.T) {
	t.Parallel()

	client := NewClient("127.0.0.1:0", "")
	defer func() { _ = client.Close() }()

	if err := client.EnqueueSendPush(PushPayload{UserID: 1, Title: "title", Body: "body"}); err == nil {
		t.Fatal("expected send push enqueue to fail without redis")
	}

	if err := client.EnqueueCheckBudget(BudgetCheckPayload{UserID: 1}); err == nil {
		t.Fatal("expected check budget enqueue to fail without redis")
	}
}
