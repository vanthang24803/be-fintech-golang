package models

import "time"

type SavingsGoal struct {
	ID           int64      `db:"id" json:"id,string"`
	UserID       int64      `db:"user_id" json:"user_id,string"`
	Name         string     `db:"name" json:"name"`
	TargetAmount float64    `db:"target_amount" json:"target_amount"`
	CurrentAmount float64   `db:"current_amount" json:"current_amount"`
	TargetDate   *time.Time `db:"target_date" json:"target_date"`
	Status       string     `db:"status" json:"status"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`

	// Derived
	ProgressPercentage float64 `json:"progress_percentage"`
}

type GoalContribution struct {
	ID        int64     `db:"id" json:"id,string"`
	GoalID    int64     `db:"goal_id" json:"goal_id,string"`
	FundID    *int64    `db:"fund_id" json:"fund_id,string,omitempty"`
	Amount    float64   `db:"amount" json:"amount"`
	Type      string    `db:"type" json:"type"` // deposit, withdrawal
	Notes     string    `db:"notes" json:"notes"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateGoalRequest struct {
	Name         string     `json:"name"          validate:"required,min=1,max=100"`
	TargetAmount float64    `json:"target_amount" validate:"required,gt=0"`
	TargetDate   *time.Time `json:"target_date"`
}

type GoalContributeRequest struct {
	GoalID int64   `json:"goal_id,string" validate:"required"`
	FundID int64   `json:"fund_id,string" validate:"required"`
	Amount float64 `json:"amount"        validate:"required,gt=0"`
	Notes  string  `json:"notes"`
}

type GoalWithdrawRequest struct {
	GoalID int64   `json:"goal_id,string" validate:"required"`
	FundID int64   `json:"fund_id,string" validate:"required"`
	Amount float64 `json:"amount"        validate:"required,gt=0"`
}

type GoalResponse struct {
	Goal          *SavingsGoal        `json:"goal"`
	Contributions []GoalContribution `json:"contributions,omitempty"`
}
