package mailer

import (
	"context"
	"log"
)

type MockMailer struct{}

func NewMockMailer() *MockMailer {
	return &MockMailer{}
}

func (m *MockMailer) Send(ctx context.Context, email Email) error {
	log.Printf("[MockMailer] Sending email to: %v", email.To)
	log.Printf("[MockMailer] Subject: %s", email.Subject)
	log.Printf("[MockMailer] Body: %s", email.Body)
	return nil
}
