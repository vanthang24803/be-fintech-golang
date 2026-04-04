package mailer

import "context"

// Email contains the data needed to send an email
type Email struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

// Mailer defines the interface for sending emails
type Mailer interface {
	Send(ctx context.Context, email Email) error
}
