package models

import (
	"encoding/json"
	"time"
)

// NotificationSource identifies which feature module emitted the notification
type NotificationSource string

const (
	SourceTransaction NotificationSource = "transaction"
	SourceFund        NotificationSource = "fund"
	SourceAuth        NotificationSource = "auth"
	SourceSystem      NotificationSource = "system"
)

// NotificationType indicates the severity / intent of the notification
type NotificationType string

const (
	NotifInfo    NotificationType = "info"
	NotifSuccess NotificationType = "success"
	NotifWarning NotificationType = "warning"
	NotifAlert   NotificationType = "alert"
)

// Notification represents a single notification entry for a user.
// Multiple feature modules can push notifications without schema changes —
// use the Source field to identify the origin and Metadata for extra data.
type Notification struct {
	ID       int64              `db:"id"        json:"id,string"`
	UserID   int64              `db:"user_id"   json:"user_id,string"`

	// Origin
	Source   NotificationSource `db:"source"    json:"source"`
	SourceID *int64             `db:"source_id" json:"source_id,omitempty"` // FK to triggering entity

	// Content
	Type     NotificationType   `db:"type"      json:"type"`
	Title    string             `db:"title"     json:"title"`
	Body     string             `db:"body"      json:"body"`
	Metadata json.RawMessage    `db:"metadata"  json:"metadata,omitempty"` // source-specific extra data

	// Read state
	IsRead  bool       `db:"is_read"  json:"is_read"`
	ReadAt  *time.Time `db:"read_at"  json:"read_at,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// CreateNotificationRequest is used internally by other services to emit a notification
type CreateNotificationRequest struct {
	UserID   int64              `json:"user_id"`
	Source   NotificationSource `json:"source"   validate:"required,oneof=transaction fund auth system"`
	SourceID *int64             `json:"source_id"`
	Type     NotificationType   `json:"type"     validate:"required,oneof=info success warning alert"`
	Title    string             `json:"title"    validate:"required,min=1,max=200"`
	Body     string             `json:"body"     validate:"required"`
	Metadata json.RawMessage    `json:"metadata"`
}

// NotificationFilter holds optional query params for listing notifications
type NotificationFilter struct {
	Source NotificationSource
	IsRead *bool // nil = all, true = read only, false = unread only
}

// MarkReadRequest is the payload for marking one or more notifications as read
type MarkReadRequest struct {
	IDs []int64 `json:"ids" validate:"required,min=1"`
}
