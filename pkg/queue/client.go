package queue

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

// Task Types
const (
	TypeSendPush   = "task:send_push"
	TypeCheckBudget = "task:check_budget"
)

// TaskPayloads
type PushPayload struct {
	UserID int64             `json:"user_id"`
	Title  string            `json:"title"`
	Body   string            `json:"body"`
	Data   map[string]string `json:"data"`
}

type BudgetCheckPayload struct {
	UserID     int64  `json:"user_id"`
	CategoryID *int64 `json:"category_id"`
}

// Client wraps asynq.Client for easy task enqueuing
type Client struct {
	client *asynq.Client
}

func NewClient(redisAddr, redisPassword string) *Client {
	return &Client{
		client: asynq.NewClient(asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPassword,
		}),
	}
}

func (c *Client) Close() error {
	return c.client.Close()
}

// EnqueueSendPush enqueues a push notification task
func (c *Client) EnqueueSendPush(payload PushPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeSendPush, data)
	_, err = c.client.Enqueue(task)
	return err
}

// EnqueueCheckBudget enqueues a budget threshold check task
func (c *Client) EnqueueCheckBudget(payload BudgetCheckPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeCheckBudget, data)
	_, err = c.client.Enqueue(task)
	return err
}
