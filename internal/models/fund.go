package models

import "time"

// Fund represents a virtual "money pocket" for a specific saving goal
type Fund struct {
	ID           int64     `db:"id"            json:"id,string"`
	UserID       int64     `db:"user_id"        json:"user_id,string"`
	Name         string    `db:"name"           json:"name"`
	Description  *string   `db:"description"    json:"description,omitempty"`
	TargetAmount float64   `db:"target_amount"  json:"target_amount"`
	Balance      float64   `db:"balance"        json:"balance"`
	Currency     string    `db:"currency"       json:"currency"`
	CreatedAt    time.Time `db:"created_at"     json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"     json:"updated_at"`
}

// CreateFundRequest is the payload for creating a new fund
type CreateFundRequest struct {
	Name         string  `json:"name"          validate:"required,min=1,max=100"`
	Description  *string `json:"description"`
	TargetAmount float64 `json:"target_amount" validate:"omitempty,min=0"` // 0 = no target
	Balance      float64 `json:"balance"       validate:"omitempty,min=0"` // initial deposit (optional)
	Currency     string  `json:"currency"      validate:"omitempty,len=3"` // default: VND
}

// UpdateFundRequest is the payload for updating fund metadata
type UpdateFundRequest struct {
	Name         string  `json:"name"          validate:"required,min=1,max=100"`
	Description  *string `json:"description"`
	TargetAmount float64 `json:"target_amount" validate:"omitempty,min=0"`
	Currency     string  `json:"currency"      validate:"omitempty,len=3"`
}

// FundTransactionRequest is used for both deposit and withdraw operations
type FundTransactionRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Note   *string `json:"note"`
}
