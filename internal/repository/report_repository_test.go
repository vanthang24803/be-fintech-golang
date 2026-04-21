package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestReportRepository_Queries(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewReportRepository(db)
	start := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery(quotedSQL(`
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
	`)).
		WithArgs(int64(42), start, end).
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "category_name", "category_icon", "total_amount"}).
			AddRow(int64(1), "Food", "utensils", 100.0))
	summary, err := repo.GetCategorySummary(context.Background(), 42, start, end)
	if err != nil || len(summary) != 1 || summary[0].CategoryName != "Food" {
		t.Fatalf("GetCategorySummary() = %+v, %v", summary, err)
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM') as month,
			SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as income,
			SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as expense
		FROM transactions
		WHERE user_id = $1 AND created_at >= $2
		GROUP BY month
		ORDER BY month DESC
	`)).
		WithArgs(int64(42), start).
		WillReturnRows(sqlmock.NewRows([]string{"month", "income", "expense"}).
			AddRow("2026-04", 500.0, 100.0))
	trend, err := repo.GetMonthlyTrend(context.Background(), 42, start)
	if err != nil || len(trend) != 1 || trend[0].Month != "2026-04" {
		t.Fatalf("GetMonthlyTrend() = %+v, %v", trend, err)
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT
			TO_CHAR(transaction_date, 'YYYY-MM-DD') as date,
			SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as income,
			SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as expense
		FROM transactions
		WHERE user_id = $1 AND transaction_date >= $2
		GROUP BY date
		ORDER BY date ASC
	`)).
		WithArgs(int64(42), start).
		WillReturnRows(sqlmock.NewRows([]string{"date", "income", "expense"}).
			AddRow("2026-04-10", 500.0, 120.0))
	dailyTrend, err := repo.GetDailyTrend(context.Background(), 42, start)
	if err != nil || len(dailyTrend) != 1 || dailyTrend[0].Income != 500 || dailyTrend[0].Expense != 120 {
		t.Fatalf("GetDailyTrend() = %+v, %v", dailyTrend, err)
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT
			COALESCE(c.id, 0) as category_id,
			COALESCE(c.name, 'Uncategorized') as category_name,
			SUM(t.amount) as amount
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = $1 AND t.type = 'income' AND t.transaction_date >= $2 AND t.transaction_date <= $3
		GROUP BY COALESCE(c.id, 0), COALESCE(c.name, 'Uncategorized')
		ORDER BY amount DESC
		 LIMIT $4`)).
		WithArgs(int64(42), start, end, 5).
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "category_name", "amount"}).
			AddRow(int64(0), "Uncategorized", 200.0))
	breakdown, err := repo.GetIncomeCategoryBreakdown(context.Background(), 42, start, end, 5)
	if err != nil || len(breakdown) != 1 || breakdown[0].Amount != 200 {
		t.Fatalf("GetIncomeCategoryBreakdown() = %+v, %v", breakdown, err)
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT
			DATE(t.transaction_date) as date,
			SUM(t.amount) as amount
		FROM transactions t
		WHERE t.user_id = $1 AND t.type = 'income' AND t.category_id = $2 AND t.transaction_date >= $3 AND t.transaction_date <= $4
		GROUP BY DATE(t.transaction_date)
		ORDER BY DATE(t.transaction_date) ASC
	`)).
		WithArgs(int64(42), int64(8), start, end).
		WillReturnRows(sqlmock.NewRows([]string{"date", "amount"}).
			AddRow("2026-04-10", 150.0))
	points, err := repo.GetCategoryTrend(context.Background(), 42, 8, start, end, "day")
	if err != nil || len(points) != 1 || points[0].Amount != 150 {
		t.Fatalf("GetCategoryTrend() = %+v, %v", points, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}
