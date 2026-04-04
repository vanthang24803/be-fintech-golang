package models

import "time"

// Transaction represents a financial transaction (income or expense)
type Transaction struct {
	ID              int64     `db:"id" json:"id,string"`
	UserID          int64     `db:"user_id" json:"user_id,string"`
	SourcePaymentID int64     `db:"sourcepayment_id" json:"source_payment_id,string"`
	CategoryID      *int64    `db:"category_id" json:"category_id,omitempty"`
	Amount          float64   `db:"amount" json:"amount"`
	Type            string    `db:"type" json:"type"` // "income" or "expense"
	Description     *string   `db:"description" json:"description,omitempty"`
	TransactionDate time.Time `db:"transaction_date" json:"transaction_date"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// TransactionDetail enriches Transaction with source and category names for display
type TransactionDetail struct {
	Transaction
	SourceName   string  `db:"source_name" json:"source_name"`
	CategoryName *string `db:"category_name" json:"category_name,omitempty"`
}

// CreateTransactionRequest is the payload for creating a new transaction
type CreateTransactionRequest struct {
	SourcePaymentID int64     `json:"source_payment_id,string" validate:"required"`
	CategoryID      *int64    `json:"category_id,omitempty"`
	Amount          float64   `json:"amount" validate:"required,gt=0"`
	Type            string    `json:"type" validate:"required,oneof=income expense"`
	Description     *string   `json:"description,omitempty"`
	TransactionDate time.Time `json:"transaction_date" validate:"required"`
}

// UpdateTransactionRequest is the payload for updating a transaction
type UpdateTransactionRequest struct {
	SourcePaymentID int64     `json:"source_payment_id,string" validate:"required"`
	CategoryID      *int64    `json:"category_id,omitempty"`
	Amount          float64   `json:"amount" validate:"required,gt=0"`
	Type            string    `json:"type" validate:"required,oneof=income expense"`
	Description     *string   `json:"description,omitempty"`
	TransactionDate time.Time `json:"transaction_date" validate:"required"`
}

// TransactionFilter holds optional query params for listing transactions
type TransactionFilter struct {
	Type            string
	CategoryID      int64
	SourcePaymentID int64
}
