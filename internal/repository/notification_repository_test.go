package repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maynguyen24/sever/internal/models"
)

func TestNotificationRepository_Operations(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewNotificationRepository(db)
	now := time.Now()
	sourceID := int64(5)
	meta := json.RawMessage(`{"k":"v"}`)
	notif := &models.Notification{
		UserID:   42,
		Source:   models.SourceSystem,
		SourceID: &sourceID,
		Type:     models.NotifInfo,
		Title:    "Hello",
		Body:     "Body",
		Metadata: meta,
	}

	mock.ExpectQuery(quotedSQL(`
		INSERT INTO notifications (id, user_id, source, source_id, type, title, body, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at
	`)).
		WithArgs(sqlmock.AnyArg(), notif.UserID, notif.Source, notif.SourceID, notif.Type, notif.Title, notif.Body, notif.Metadata).
		WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(now))
	if err := repo.Create(context.Background(), notif); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	isRead := false
	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, source, source_id, type, title, body, metadata, is_read, read_at, created_at 
		FROM notifications WHERE user_id = $1 AND source = $2 AND is_read = $3 ORDER BY created_at DESC`)).
		WithArgs(int64(42), models.SourceSystem, isRead).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "source", "source_id", "type", "title", "body", "metadata", "is_read", "read_at", "created_at"}).
			AddRow(int64(1), int64(42), models.SourceSystem, sourceID, models.NotifInfo, "Hello", "Body", meta, false, nil, now))
	list, err := repo.GetByUserID(context.Background(), 42, models.NotificationFilter{Source: models.SourceSystem, IsRead: &isRead})
	if err != nil || len(list) != 1 {
		t.Fatalf("GetByUserID() = %+v, %v", list, err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	count, err := repo.GetUnreadCount(context.Background(), 42)
	if err != nil || count != 3 {
		t.Fatalf("GetUnreadCount() = %d, %v", count, err)
	}

	mock.ExpectExec(`UPDATE notifications SET is_read = TRUE, read_at = NOW\(\) 
		WHERE user_id = \? AND id IN \(\?, \?\) AND is_read = FALSE`).
		WithArgs(int64(42), int64(1), int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 2))
	if err := repo.MarkAsRead(context.Background(), 42, []int64{1, 2}); err != nil {
		t.Fatalf("MarkAsRead() error = %v", err)
	}

	mock.ExpectExec(quotedSQL(`DELETE FROM notifications WHERE id = $1 AND user_id = $2`)).
		WithArgs(int64(1), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.Delete(context.Background(), 42, 1); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}
