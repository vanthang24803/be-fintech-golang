package push

import (
	"context"
	"testing"
)

func TestNewFirebaseClient_RequiresCredentialsPath(t *testing.T) {
	t.Parallel()

	if _, err := NewFirebaseClient(""); err == nil {
		t.Fatal("expected error for empty credentials path")
	}
}

func TestMockPushClient_NoOp(t *testing.T) {
	t.Parallel()

	client := &MockPushClient{}
	if err := client.SendPush(context.Background(), "token", "title", "body", map[string]string{"k": "v"}); err != nil {
		t.Fatalf("SendPush() error = %v", err)
	}
	if err := client.SendToTopic(context.Background(), "topic", "title", "body", nil); err != nil {
		t.Fatalf("SendToTopic() error = %v", err)
	}
}
