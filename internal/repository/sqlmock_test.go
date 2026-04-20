package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New(): %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	return sqlx.NewDb(db, "sqlmock"), mock
}

func quotedSQL(query string) string {
	return regexp.QuoteMeta(query)
}

var (
	userCols         = []string{"id", "username", "email", "password_hash", "google_id", "created_at", "updated_at"}
	profileCols      = []string{"id", "user_id", "full_name", "avatar_url", "phone_number", "date_of_birth", "created_at", "updated_at"}
	fundCols         = []string{"id", "user_id", "name", "description", "target_amount", "balance", "currency", "created_at", "updated_at"}
	sourceCols       = []string{"id", "user_id", "name", "type", "balance", "currency", "created_at", "updated_at"}
	categoryCols     = []string{"id", "user_id", "name", "type", "icon", "created_at", "updated_at"}
	budgetCols       = []string{"id", "user_id", "category_id", "amount", "period", "start_date", "end_date", "is_active", "created_at", "updated_at"}
	deviceCols       = []string{"id", "user_id", "device_fingerprint", "device_name", "platform", "push_token", "fido_credential_id", "fido_public_key", "fido_sign_count", "fido_aaguid", "is_trusted", "is_active", "last_used_at", "created_at", "updated_at"}
	transactionCols  = []string{"id", "user_id", "sourcepayment_id", "category_id", "amount", "type", "description", "transaction_date", "created_at", "updated_at"}
	transactionDCols = []string{"id", "user_id", "sourcepayment_id", "category_id", "amount", "type", "description", "transaction_date", "created_at", "updated_at", "source_name", "category_name"}
)

func timestampRows(now time.Time) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"created_at", "updated_at"}).AddRow(now, now)
}
