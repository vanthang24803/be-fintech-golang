package models

import (
	"time"
)

// ReportRequest defines common filters for report queries
type ReportRequest struct {
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date"   validate:"required,gtefield=StartDate"`
}

// MonthlyTrendRequest defines the number of months to look back
type MonthlyTrendRequest struct {
	Months int `json:"months" validate:"omitempty,min=1,max=24"` // Default: 6 if zero
}

// CategorySummary provides aggregated spending data for a category
type CategorySummary struct {
	CategoryID   int64   `json:"category_id" db:"category_id"`
	CategoryName string  `json:"category_name" db:"category_name"`
	CategoryIcon string  `json:"category_icon" db:"category_icon"`
	TotalAmount  float64 `json:"total_amount" db:"total_amount"`
	Percentage   float64 `json:"percentage"`
}

// MonthlySummary provides income/expense totals for a specific month
type MonthlySummary struct {
	Month     string  `json:"month" db:"month"` // Format: YYYY-MM
	Income    float64 `json:"income" db:"income"`
	Expense   float64 `json:"expense" db:"expense"`
	NetProfit float64 `json:"net_profit"`
}

type IncomeCategoryBreakdownRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Limit     int       `json:"limit"`
}

type IncomeCategoryBreakdownItem struct {
	CategoryID   int64   `json:"category_id" db:"category_id"`
	CategoryName string  `json:"category_name" db:"category_name"`
	Amount       float64 `json:"amount" db:"amount"`
	Percentage   float64 `json:"percentage"`
}

type IncomeCategoryBreakdownResponse struct {
	TotalIncome float64                        `json:"total_income"`
	RangeStart  time.Time                      `json:"range_start"`
	RangeEnd    time.Time                      `json:"range_end"`
	Items       []*IncomeCategoryBreakdownItem `json:"items"`
}

type CategoryTrendRequest struct {
	CategoryID  int64     `json:"category_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Granularity string    `json:"granularity"`
}

type CategoryTrendPoint struct {
	Date   string  `json:"date" db:"date"`
	Amount float64 `json:"amount" db:"amount"`
}

type CategoryTrendResponse struct {
	CategoryID   int64                 `json:"category_id"`
	CategoryName string                `json:"category_name"`
	Granularity  string                `json:"granularity"`
	Points       []*CategoryTrendPoint `json:"points"`
}
