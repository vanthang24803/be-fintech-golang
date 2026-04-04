package mailer

import (
	"context"
	"fmt"
	"net/smtp"
)

type SMTPMailer struct {
	host     string
	port     string
	user     string
	password string
	from     string
}

// NewSMTPMailer initializes a new SMTP mailer
func NewSMTPMailer(host, port, user, password, from string) *SMTPMailer {
	return &SMTPMailer{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
	}
}

func (m *SMTPMailer) Send(ctx context.Context, email Email) error {
	// 1. Build message
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	if email.IsHTML {
		mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	}

	subject := "Subject: " + email.Subject + "\n"
	msg := []byte(subject + mime + email.Body)

	// 2. Setup auth
	auth := smtp.PlainAuth("", m.user, m.password, m.host)

	// 3. Send email
	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	err := smtp.SendMail(addr, auth, m.from, email.To, msg)
	if err != nil {
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}

	return nil
}
