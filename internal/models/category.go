package models

import "time"

// Category represents a transaction category (income/expense)
type Category struct {
	ID        int64     `db:"id" json:"id,string"`
	UserID    *int64    `db:"user_id" json:"user_id,omitempty"` // nil = system default category
	Name      string    `db:"name" json:"name"`
	Type      string    `db:"type" json:"type"` // "income" or "expense"
	Icon      *string   `db:"icon" json:"icon,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// CreateCategoryRequest is the payload for creating a new category
type CreateCategoryRequest struct {
	Name string  `json:"name" validate:"required,min=1,max=100"`
	Type string  `json:"type" validate:"required,oneof=income expense"`
	Icon *string `json:"icon"`
}

// UpdateCategoryRequest is the payload for updating a category
type UpdateCategoryRequest struct {
	Name string  `json:"name" validate:"required,min=1,max=100"`
	Type string  `json:"type" validate:"required,oneof=income expense"`
	Icon *string `json:"icon"`
}
