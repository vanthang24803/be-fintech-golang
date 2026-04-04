package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/internal/models"
)

type SavingsGoalRepository struct {
	db *sqlx.DB
}

func NewSavingsGoalRepository(db *sqlx.DB) *SavingsGoalRepository {
	return &SavingsGoalRepository{db: db}
}

func (r *SavingsGoalRepository) CreateGoal(ctx context.Context, goal *models.SavingsGoal) error {
	query := `
		INSERT INTO savings_goals (id, user_id, name, target_amount, target_date, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, query, 
		goal.ID, goal.UserID, goal.Name, goal.TargetAmount, goal.TargetDate, goal.Status,
	).Scan(&goal.CreatedAt, &goal.UpdatedAt)
}

func (r *SavingsGoalRepository) GetGoalByID(ctx context.Context, id int64) (*models.SavingsGoal, error) {
	var goal models.SavingsGoal
	query := `SELECT id, user_id, name, target_amount, current_amount, target_date, status, created_at, updated_at FROM savings_goals WHERE id = $1`
	err := r.db.GetContext(ctx, &goal, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &goal, nil
}

func (r *SavingsGoalRepository) ListGoals(ctx context.Context, userID int64) ([]models.SavingsGoal, error) {
	var goals []models.SavingsGoal
	query := `SELECT id, user_id, name, target_amount, current_amount, target_date, status, created_at, updated_at FROM savings_goals WHERE user_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &goals, query, userID)
	return goals, err
}

func (r *SavingsGoalRepository) UpdateGoalAmount(ctx context.Context, goalID int64, amount float64) error {
	query := `UPDATE savings_goals SET current_amount = current_amount + $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, amount, goalID)
	return err
}

func (r *SavingsGoalRepository) CreateContribution(ctx context.Context, c *models.GoalContribution) error {
	query := `INSERT INTO goal_contributions (id, goal_id, fund_id, amount, type, notes) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, c.ID, c.GoalID, c.FundID, c.Amount, c.Type, c.Notes)
	return err
}

func (r *SavingsGoalRepository) GetContributionsByGoal(ctx context.Context, goalID int64) ([]models.GoalContribution, error) {
	var contributions []models.GoalContribution
	query := `SELECT id, goal_id, fund_id, amount, type, notes, created_at FROM goal_contributions WHERE goal_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &contributions, query, goalID)
	return contributions, err
}

func (r *SavingsGoalRepository) DeleteGoal(ctx context.Context, id int64) error {
	query := "DELETE FROM savings_goals WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
