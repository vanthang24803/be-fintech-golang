package repository

import (
	"context"
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
func (r *ReportRepository) GetCategorySummary(ctx context.Context, userID int64, start, end time.Time) ([]*models.CategorySummary, error) {
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
	if err := r.db.SelectContext(ctx, &summary, query, userID, start, end); err != nil {
		return nil, err
	}
	return summary, nil
}

// GetMonthlyTrend aggregates income and expense totals grouped by month
func (r *ReportRepository) GetMonthlyTrend(ctx context.Context, userID int64, since time.Time) ([]*models.MonthlySummary, error) {
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
	if err := r.db.SelectContext(ctx, &trend, query, userID, since); err != nil {
		return nil, err
	}
	return trend, nil
}

// GetIncomeCategoryBreakdown aggregates income by category for a period
func (r *ReportRepository) GetIncomeCategoryBreakdown(ctx context.Context, userID int64, start, end time.Time, limit int) ([]*models.IncomeCategoryBreakdownItem, error) {
	var summary []*models.IncomeCategoryBreakdownItem
	query := `
		SELECT
			COALESCE(c.id, 0) as category_id,
			COALESCE(c.name, 'Uncategorized') as category_name,
			SUM(t.amount) as amount
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = $1 AND t.type = 'income' AND t.transaction_date >= $2 AND t.transaction_date <= $3
		GROUP BY COALESCE(c.id, 0), COALESCE(c.name, 'Uncategorized')
		ORDER BY amount DESC
	`
	args := []any{userID, start, end}
	if limit > 0 {
		query += ` LIMIT $4`
		args = append(args, limit)
	}
	if err := r.db.SelectContext(ctx, &summary, query, args...); err != nil {
		return nil, err
	}
	return summary, nil
}

// GetDailyTrend aggregates income and expense totals grouped by day
func (r *ReportRepository) GetDailyTrend(ctx context.Context, userID int64, since time.Time) ([]*models.DailySummary, error) {
	var trend []*models.DailySummary
	query := `
		SELECT
			TO_CHAR(transaction_date, 'YYYY-MM-DD') as date,
			SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as income,
			SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as expense
		FROM transactions
		WHERE user_id = $1 AND transaction_date >= $2
		GROUP BY date
		ORDER BY date ASC
	`
	if err := r.db.SelectContext(ctx, &trend, query, userID, since); err != nil {
		return nil, err
	}
	return trend, nil
}

// GetCategoryTrend aggregates income totals grouped by day for a category
func (r *ReportRepository) GetCategoryTrend(ctx context.Context, userID, categoryID int64, start, end time.Time, granularity string) ([]*models.CategoryTrendPoint, error) {
	var points []*models.CategoryTrendPoint
	query := `
		SELECT
			DATE(t.transaction_date) as date,
			SUM(t.amount) as amount
		FROM transactions t
		WHERE t.user_id = $1 AND t.type = 'income' AND t.category_id = $2 AND t.transaction_date >= $3 AND t.transaction_date <= $4
		GROUP BY DATE(t.transaction_date)
		ORDER BY DATE(t.transaction_date) ASC
	`
	if err := r.db.SelectContext(ctx, &points, query, userID, categoryID, start, end); err != nil {
		return nil, err
	}
	return points, nil
}
