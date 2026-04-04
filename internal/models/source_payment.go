package models

import "time"

// SourcePayment represents a payment source (wallet or bank account)
type SourcePayment struct {
	ID        int64     `db:"id" json:"id,string"`
	UserID    int64     `db:"user_id" json:"user_id,string"`
	Name      string    `db:"name" json:"name"`
	Type      string    `db:"type" json:"type"`
	Balance   float64   `db:"balance" json:"balance"`
	Currency  string    `db:"currency" json:"currency"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// CreateSourcePaymentRequest is the payload for creating a new source
type CreateSourcePaymentRequest struct {
	Name     string  `json:"name"     validate:"required,min=1,max=100"`
	Type     string  `json:"type"     validate:"required"` // e.g. "wallet", "bank", "credit_card"
	Balance  float64 `json:"balance"  validate:"omitempty,min=0"`
	Currency string  `json:"currency" validate:"omitempty,len=3"` // Default: VND
}

// UpdateSourcePaymentRequest is the payload for updating a source
type UpdateSourcePaymentRequest struct {
	Name     string  `json:"name"     validate:"required,min=1,max=100"`
	Type     string  `json:"type"     validate:"required"`
	Currency string  `json:"currency" validate:"omitempty,len=3"`
}
