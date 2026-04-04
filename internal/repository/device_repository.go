package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
	"github.com/maynguyen24/sever/pkg/snowflake"
)

// DeviceRepository handles database operations for devices
type DeviceRepository struct {
	db *sqlx.DB
}

func NewDeviceRepository(db *sqlx.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// Create inserts a new device record into the database
func (r *DeviceRepository) Create(ctx context.Context, device *models.Device) error {
	device.ID = snowflake.GenerateID()

	query := `
		INSERT INTO devices (id, user_id, device_fingerprint, device_name, platform, push_token)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, query,
		device.ID, device.UserID, device.DeviceFingerprint, device.DeviceName,
		device.Platform, device.PushToken,
	).Scan(&device.CreatedAt, &device.UpdatedAt)
}

// GetByFingerprint fetches a device record by its unique fingerprint
func (r *DeviceRepository) GetByFingerprint(ctx context.Context, fingerprint string) (*models.Device, error) {
	var device models.Device
	query := `SELECT id, user_id, device_fingerprint, device_name, platform, push_token, 
		fido_credential_id, fido_public_key, fido_sign_count, fido_aaguid, is_trusted, is_active, last_used_at, created_at, updated_at
		FROM devices WHERE device_fingerprint = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &device, query, fingerprint)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

// GetByUserID fetches all devices for a specific user
func (r *DeviceRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Device, error) {
	var devices []*models.Device
	query := `SELECT id, user_id, device_fingerprint, device_name, platform, push_token,
		fido_credential_id, fido_public_key, fido_sign_count, fido_aaguid, is_trusted, is_active, last_used_at, created_at, updated_at
		FROM devices WHERE user_id = $1 ORDER BY last_used_at DESC NULLS LAST, created_at DESC`
	if err := r.db.SelectContext(ctx, &devices, query, userID); err != nil {
		return nil, err
	}
	return devices, nil
}

// GetPushTokensByUserID retrieves all unique push tokens for a user
func (r *DeviceRepository) GetPushTokensByUserID(ctx context.Context, userID int64) ([]string, error) {
	var tokens []string
	query := `SELECT DISTINCT push_token FROM devices WHERE user_id = $1 AND push_token IS NOT NULL AND push_token != ''`
	err := r.db.SelectContext(ctx, &tokens, query, userID)
	return tokens, err
}

// GetByID fetches a specific device record for a user
func (r *DeviceRepository) GetByID(ctx context.Context, id, userID int64) (*models.Device, error) {
	var device models.Device
	query := `SELECT id, user_id, device_fingerprint, device_name, platform, push_token,
		fido_credential_id, fido_public_key, fido_sign_count, fido_aaguid, is_trusted, is_active, last_used_at, created_at, updated_at
		FROM devices WHERE id = $1 AND user_id = $2 LIMIT 1`
	err := r.db.GetContext(ctx, &device, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

// Delete removes a device record
func (r *DeviceRepository) Delete(ctx context.Context, id, userID int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM devices WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

// UpdateLastUsed updates the last_used_at timestamp for a device
func (r *DeviceRepository) UpdateLastUsed(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE devices SET last_used_at = NOW() WHERE id = $1`, id)
	return err
}

// GetByCredentialID fetches the device that owns a given FIDO2 credential ID
func (r *DeviceRepository) GetByCredentialID(ctx context.Context, credentialID string) (*models.Device, error) {
	var device models.Device
	query := `SELECT id, user_id, device_fingerprint, device_name, platform, push_token,
		fido_credential_id, fido_public_key, fido_sign_count, fido_aaguid, is_trusted, is_active, last_used_at, created_at, updated_at
		FROM devices WHERE fido_credential_id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &device, query, credentialID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

// UpdateFIDOCredential stores a newly enrolled FIDO2 credential on a device
func (r *DeviceRepository) UpdateFIDOCredential(ctx context.Context, deviceID int64, credentialID, publicKey, aaguid string, signCount int64) error {
	_, err := r.db.ExecContext(ctx, 
		`UPDATE devices SET fido_credential_id = $1, fido_public_key = $2, fido_aaguid = $3, fido_sign_count = $4, updated_at = NOW()
		 WHERE id = $5`,
		credentialID, publicKey, aaguid, signCount, deviceID,
	)
	return err
}

// UpdateSignCount updates the FIDO2 authenticator sign count after a successful assertion
func (r *DeviceRepository) UpdateSignCount(ctx context.Context, deviceID int64, signCount int64) error {
	_, err := r.db.ExecContext(ctx, 
		`UPDATE devices SET fido_sign_count = $1, updated_at = NOW() WHERE id = $2`,
		signCount, deviceID,
	)
	return err
}
