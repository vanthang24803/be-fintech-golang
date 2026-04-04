package push

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// PushClient defines the interface for sending push notifications
type PushClient interface {
	SendPush(ctx context.Context, token, title, body string, data map[string]string) error
	SendToTopic(ctx context.Context, topic, title, body string, data map[string]string) error
}

// FirebaseClient implements PushClient using Firebase Cloud Messaging
type FirebaseClient struct {
	app       *firebase.App
	messaging *messaging.Client
}

// NewFirebaseClient initializes a new Firebase messaging client
func NewFirebaseClient(credentialsPath string) (*FirebaseClient, error) {
	if credentialsPath == "" {
		return nil, fmt.Errorf("Firebase credentials path is empty")
	}

	opt := option.WithServiceAccountFile(credentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %w", err)
	}

	msgClient, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting messaging client: %w", err)
	}

	return &FirebaseClient{
		app:       app,
		messaging: msgClient,
	}, nil
}

// SendPush sends a notification to a specific device token
func (c *FirebaseClient) SendPush(ctx context.Context, token, title, body string, data map[string]string) error {
	if token == "" {
		return nil // Ignore empty tokens
	}

	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	response, err := c.messaging.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("error sending fcm message: %w", err)
	}

	log.Printf("Successfully sent fcm message: %s", response)
	return nil
}

// SendToTopic sends a notification to all devices subscribed to a topic
func (c *FirebaseClient) SendToTopic(ctx context.Context, topic, title, body string, data map[string]string) error {
	message := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	response, err := c.messaging.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("error sending fcm topic message: %w", err)
	}

	log.Printf("Successfully sent fcm topic message: %s", response)
	return nil
}

// MockPushClient is used for testing or when firebase is disabled
type MockPushClient struct{}

func (m *MockPushClient) SendPush(ctx context.Context, token, title, body string, data map[string]string) error {
	log.Printf("[MOCK PUSH] To: %s, Title: %s, Body: %s", token, title, body)
	return nil
}

func (m *MockPushClient) SendToTopic(ctx context.Context, topic, title, body string, data map[string]string) error {
	log.Printf("[MOCK TOPIC] Topic: %s, Title: %s, Body: %s", topic, title, body)
	return nil
}
