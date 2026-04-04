package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/internal/models"
)

// ReportRepository handles analytical database queries
type ReportRepository struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetCategorySummary aggregates expenses by category for a period
func (r *ReportRepository) GetCategorySummary(userID int64, start, end time.Time) ([]*models.CategorySummary, error) {
	var summary []*models.CategorySummary
	query := `
		SELECT 
			c.id as category_id, 
			c.name as category_name, 
			c.icon as category_icon, 
			SUM(t.amount) as total_amount
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = $1 AND t.type = 'expense' AND t.created_at >= $2 AND t.created_at <= $3
		GROUP BY c.id, c.name, c.icon
		ORDER BY total_amount DESC
	`
	if err := r.db.Select(&summary, query, userID, start, end); err != nil {
		return nil, err
	}
	return summary, nil
}

// GetMonthlyTrend aggregates income and expense totals grouped by month
func (r *ReportRepository) GetMonthlyTrend(userID int64, since time.Time) ([]*models.MonthlySummary, error) {
	var trend []*models.MonthlySummary
	query := `
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM') as month,
			SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as income,
			SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as expense
		FROM transactions
		WHERE user_id = $1 AND created_at >= $2
		GROUP BY month
		ORDER BY month DESC
	`
	if err := r.db.Select(&trend, query, userID, since); err != nil {
		return nil, err
	}
	return trend, nil
}
