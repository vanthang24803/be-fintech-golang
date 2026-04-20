package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maynguyen24/sever/internal/models"
)

func TestSavingsGoalRepository_Operations(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewSavingsGoalRepository(db)
	now := time.Now()
	targetDate := now.AddDate(0, 6, 0)
	goal := &models.SavingsGoal{
		ID:           1,
		UserID:       42,
		Name:         "Car",
		TargetAmount: 5000,
		TargetDate:   &targetDate,
		Status:       "active",
	}
	fundID := int64(2)
	contribution := models.GoalContribution{ID: 1, GoalID: 1, FundID: &fundID, Amount: 500, Type: "deposit", Notes: "seed", CreatedAt: now}

	mock.ExpectQuery(quotedSQL(`
		INSERT INTO savings_goals (id, user_id, name, target_amount, target_date, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`)).
		WithArgs(goal.ID, goal.UserID, goal.Name, goal.TargetAmount, goal.TargetDate, goal.Status).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).AddRow(now, now))
	if err := repo.CreateGoal(context.Background(), goal); err != nil {
		t.Fatalf("CreateGoal() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, name, target_amount, current_amount, target_date, status, created_at, updated_at FROM savings_goals WHERE id = $1`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "target_amount", "current_amount", "target_date", "status", "created_at", "updated_at"}).
			AddRow(int64(1), int64(42), "Car", 5000.0, 500.0, targetDate, "active", now, now))
	got, err := repo.GetGoalByID(context.Background(), 1)
	if err != nil || got == nil || got.Name != "Car" {
		t.Fatalf("GetGoalByID() = %+v, %v", got, err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, name, target_amount, current_amount, target_date, status, created_at, updated_at FROM savings_goals WHERE user_id = $1 ORDER BY created_at DESC`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "target_amount", "current_amount", "target_date", "status", "created_at", "updated_at"}).
			AddRow(int64(1), int64(42), "Car", 5000.0, 500.0, targetDate, "active", now, now))
	list, err := repo.ListGoals(context.Background(), 42)
	if err != nil || len(list) != 1 {
		t.Fatalf("ListGoals() = %+v, %v", list, err)
	}

	mock.ExpectExec(quotedSQL(`UPDATE savings_goals SET current_amount = current_amount + $1, updated_at = NOW() WHERE id = $2`)).
		WithArgs(200.0, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.UpdateGoalAmount(context.Background(), 1, 200); err != nil {
		t.Fatalf("UpdateGoalAmount() error = %v", err)
	}

	mock.ExpectExec(quotedSQL(`INSERT INTO goal_contributions (id, goal_id, fund_id, amount, type, notes) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(contribution.ID, contribution.GoalID, contribution.FundID, contribution.Amount, contribution.Type, contribution.Notes).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.CreateContribution(context.Background(), &contribution); err != nil {
		t.Fatalf("CreateContribution() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, goal_id, fund_id, amount, type, notes, created_at FROM goal_contributions WHERE goal_id = $1 ORDER BY created_at DESC`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "goal_id", "fund_id", "amount", "type", "notes", "created_at"}).
			AddRow(contribution.ID, contribution.GoalID, fundID, contribution.Amount, contribution.Type, contribution.Notes, contribution.CreatedAt))
	contributions, err := repo.GetContributionsByGoal(context.Background(), 1)
	if err != nil || len(contributions) != 1 {
		t.Fatalf("GetContributionsByGoal() = %+v, %v", contributions, err)
	}

	mock.ExpectExec(quotedSQL(`DELETE FROM savings_goals WHERE id = $1`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.DeleteGoal(context.Background(), 1); err != nil {
		t.Fatalf("DeleteGoal() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, name, target_amount, current_amount, target_date, status, created_at, updated_at FROM savings_goals WHERE id = $1`)).
		WithArgs(int64(9)).
		WillReturnError(sql.ErrNoRows)
	got, err = repo.GetGoalByID(context.Background(), 9)
	if err != nil || got != nil {
		t.Fatalf("expected nil on missing goal, got %+v err=%v", got, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}
