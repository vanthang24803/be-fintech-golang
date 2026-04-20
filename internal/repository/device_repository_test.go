package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

func TestDeviceRepository_Operations(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewDeviceRepository(db)
	now := time.Now()
	deviceName := "iPhone"
	pushToken := "push-token"
	credentialID := "cred"
	publicKey := "pk"
	aaguid := "aaguid"
	device := &models.Device{
		UserID:            42,
		DeviceFingerprint: "fingerprint-1",
		DeviceName:        &deviceName,
		Platform:          models.PlatformIOS,
		PushToken:         &pushToken,
	}

	mock.ExpectQuery(quotedSQL(`
		INSERT INTO devices (id, user_id, device_fingerprint, device_name, platform, push_token)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`)).
		WithArgs(sqlmock.AnyArg(), device.UserID, device.DeviceFingerprint, device.DeviceName, device.Platform, device.PushToken).
		WillReturnRows(timestampRows(now))
	if err := repo.Create(context.Background(), device); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, device_fingerprint, device_name, platform, push_token, 
		fido_credential_id, fido_public_key, fido_sign_count, fido_aaguid, is_trusted, is_active, last_used_at, created_at, updated_at
		FROM devices WHERE device_fingerprint = $1 LIMIT 1`)).
		WithArgs("fingerprint-1").
		WillReturnRows(sqlmock.NewRows(deviceCols).
			AddRow(int64(1), int64(42), "fingerprint-1", deviceName, "ios", pushToken, nil, nil, int64(0), nil, true, true, now, now, now))
	got, err := repo.GetByFingerprint(context.Background(), "fingerprint-1")
	if err != nil || got == nil || got.DeviceFingerprint != "fingerprint-1" {
		t.Fatalf("GetByFingerprint() = %+v, %v", got, err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, device_fingerprint, device_name, platform, push_token,
		fido_credential_id, fido_public_key, fido_sign_count, fido_aaguid, is_trusted, is_active, last_used_at, created_at, updated_at
		FROM devices WHERE user_id = $1 ORDER BY last_used_at DESC NULLS LAST, created_at DESC`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows(deviceCols).
			AddRow(int64(1), int64(42), "fingerprint-1", deviceName, "ios", pushToken, nil, nil, int64(0), nil, true, true, now, now, now))
	list, err := repo.GetByUserID(context.Background(), 42)
	if err != nil || len(list) != 1 {
		t.Fatalf("GetByUserID() = %+v, %v", list, err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT DISTINCT push_token FROM devices WHERE user_id = $1 AND push_token IS NOT NULL AND push_token != ''`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"push_token"}).AddRow(pushToken))
	tokens, err := repo.GetPushTokensByUserID(context.Background(), 42)
	if err != nil || len(tokens) != 1 || tokens[0] != pushToken {
		t.Fatalf("GetPushTokensByUserID() = %+v, %v", tokens, err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, device_fingerprint, device_name, platform, push_token,
		fido_credential_id, fido_public_key, fido_sign_count, fido_aaguid, is_trusted, is_active, last_used_at, created_at, updated_at
		FROM devices WHERE id = $1 AND user_id = $2 LIMIT 1`)).
		WithArgs(int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows(deviceCols).
			AddRow(int64(1), int64(42), "fingerprint-1", deviceName, "ios", pushToken, credentialID, publicKey, int64(2), aaguid, true, true, now, now, now))
	got, err = repo.GetByID(context.Background(), 1, 42)
	if err != nil || got == nil || got.ID != 1 {
		t.Fatalf("GetByID() = %+v, %v", got, err)
	}

	mock.ExpectExec(quotedSQL(`UPDATE devices SET last_used_at = NOW() WHERE id = $1`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.UpdateLastUsed(context.Background(), 1); err != nil {
		t.Fatalf("UpdateLastUsed() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, device_fingerprint, device_name, platform, push_token,
		fido_credential_id, fido_public_key, fido_sign_count, fido_aaguid, is_trusted, is_active, last_used_at, created_at, updated_at
		FROM devices WHERE fido_credential_id = $1 LIMIT 1`)).
		WithArgs(credentialID).
		WillReturnRows(sqlmock.NewRows(deviceCols).
			AddRow(int64(1), int64(42), "fingerprint-1", deviceName, "ios", pushToken, credentialID, publicKey, int64(2), aaguid, true, true, now, now, now))
	got, err = repo.GetByCredentialID(context.Background(), credentialID)
	if err != nil || got == nil || got.FIDOCredentialID == nil || *got.FIDOCredentialID != credentialID {
		t.Fatalf("GetByCredentialID() = %+v, %v", got, err)
	}

	mock.ExpectExec(quotedSQL(`UPDATE devices SET fido_credential_id = $1, fido_public_key = $2, fido_aaguid = $3, fido_sign_count = $4, updated_at = NOW()
		 WHERE id = $5`)).
		WithArgs("cred-2", "pk-2", "aaguid-2", int64(3), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.UpdateFIDOCredential(context.Background(), 1, "cred-2", "pk-2", "aaguid-2", 3); err != nil {
		t.Fatalf("UpdateFIDOCredential() error = %v", err)
	}

	mock.ExpectExec(quotedSQL(`UPDATE devices SET fido_sign_count = $1, updated_at = NOW() WHERE id = $2`)).
		WithArgs(int64(4), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.UpdateSignCount(context.Background(), 1, 4); err != nil {
		t.Fatalf("UpdateSignCount() error = %v", err)
	}

	mock.ExpectExec(quotedSQL(`DELETE FROM devices WHERE id = $1 AND user_id = $2`)).
		WithArgs(int64(1), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.Delete(context.Background(), 1, 42); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}

func TestDeviceRepository_ErrorBranches(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewDeviceRepository(db)

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, device_fingerprint, device_name, platform, push_token, 
		fido_credential_id, fido_public_key, fido_sign_count, fido_aaguid, is_trusted, is_active, last_used_at, created_at, updated_at
		FROM devices WHERE device_fingerprint = $1 LIMIT 1`)).
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)
	got, err := repo.GetByFingerprint(context.Background(), "missing")
	if err != nil || got != nil {
		t.Fatalf("expected nil on missing fingerprint, got %+v err=%v", got, err)
	}

	mock.ExpectExec(quotedSQL(`DELETE FROM devices WHERE id = $1 AND user_id = $2`)).
		WithArgs(int64(9), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	if err := repo.Delete(context.Background(), 9, 42); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}
