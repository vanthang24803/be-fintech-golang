package models

import (
	"time"
)

// Budget represents a spending limit for a user/category/period
type Budget struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	CategoryID *int64    `json:"category_id,omitempty" db:"category_id"` // NULL means total budget
	Amount     float64   `json:"amount" db:"amount"`
	Period     string    `json:"period" db:"period"` // 'monthly', 'weekly', 'custom'
	StartDate  time.Time `json:"start_date" db:"start_date"`
	EndDate    time.Time `json:"end_date" db:"end_date"`
	IsActive   bool      `json:"is_active" db:"is_active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CreateBudgetRequest is used when defining a new budget
type CreateBudgetRequest struct {
	CategoryID *int64  `json:"category_id"` // Optional
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	Period     string  `json:"period" validate:"required,oneof=monthly weekly custom"` // 'monthly', 'weekly'
}

// UpdateBudgetRequest is used when modifying a budget amount
type UpdateBudgetRequest struct {
	Amount   *float64 `json:"amount" validate:"omitempty,gt=0"`
	IsActive *bool    `json:"is_active"`
}

// BudgetResponse provides details about a budget including current spending status
type BudgetResponse struct {
	Budget
	CategoryName    string  `json:"category_name,omitempty"`
	CurrentSpending float64 `json:"current_spending"`
	RemainingAmount float64 `json:"remaining_amount"`
	ProgressPercent float64 `json:"progress_percent"`
}
