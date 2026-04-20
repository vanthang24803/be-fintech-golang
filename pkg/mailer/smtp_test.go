package mailer

import (
	"context"
	"strings"
	"testing"
)

func TestNewSMTPMailer(t *testing.T) {
	t.Parallel()

	mailer := NewSMTPMailer("smtp.example.com", "587", "user", "pass", "from@example.com")
	if mailer.host != "smtp.example.com" || mailer.port != "587" || mailer.from != "from@example.com" {
		t.Fatalf("unexpected SMTP mailer config: %+v", mailer)
	}
}

func TestSMTPMailerSend_WrapsSendError(t *testing.T) {
	t.Parallel()

	mailer := NewSMTPMailer("127.0.0.1", "1", "user", "pass", "from@example.com")
	err := mailer.Send(context.Background(), Email{
		To:      []string{"to@example.com"},
		Subject: "Subject",
		Body:    "<p>Body</p>",
		IsHTML:  true,
	})
	if err == nil || !strings.Contains(err.Error(), "failed to send email via SMTP") {
		t.Fatalf("expected wrapped send error, got %v", err)
	}
}
