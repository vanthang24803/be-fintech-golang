package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/snowflake"
)

// NotificationRepository handles database operations for notifications
type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create inserts a new notification record into the database
func (r *NotificationRepository) Create(notif *models.Notification) error {
	notif.ID = snowflake.GenerateID()

	query := `
		INSERT INTO notifications (id, user_id, source, source_id, type, title, body, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at
	`
	return r.db.QueryRowx(query,
		notif.ID, notif.UserID, notif.Source, notif.SourceID,
		notif.Type, notif.Title, notif.Body, notif.Metadata,
	).Scan(&notif.CreatedAt)
}

// GetByUserID fetches all notifications for a specific user with filters
func (r *NotificationRepository) GetByUserID(userID int64, filter models.NotificationFilter) ([]*models.Notification, error) {
	query := `SELECT id, user_id, source, source_id, type, title, body, metadata, is_read, read_at, created_at 
		FROM notifications WHERE user_id = $1`
	args := []interface{}{userID}
	argIdx := 2

	if filter.Source != "" {
		query += fmt.Sprintf(" AND source = $%d", argIdx)
		args = append(args, filter.Source)
		argIdx++
	}

	if filter.IsRead != nil {
		query += fmt.Sprintf(" AND is_read = $%d", argIdx)
		args = append(args, *filter.IsRead)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	var notifications []*models.Notification
	if err := r.db.Select(&notifications, query, args...); err != nil {
		return nil, err
	}
	return notifications, nil
}

// GetUnreadCount counts unread notifications for a user
func (r *NotificationRepository) GetUnreadCount(userID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`
	if err := r.db.Get(&count, query, userID); err != nil {
		return 0, err
	}
	return count, nil
}

// MarkAsRead updates multiple notifications as read for a user
func (r *NotificationRepository) MarkAsRead(userID int64, ids []int64) error {
	query, args, err := sqlx.In(`UPDATE notifications SET is_read = TRUE, read_at = NOW() 
		WHERE user_id = ? AND id IN (?) AND is_read = FALSE`, userID, ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

// Delete removes a specific notification record
func (r *NotificationRepository) Delete(userID int64, id int64) error {
	result, err := r.db.Exec(`DELETE FROM notifications WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("notification not found")
	}
	return nil
}
