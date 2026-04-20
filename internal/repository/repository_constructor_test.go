package repository

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestRepositoryConstructors(t *testing.T) {
	t.Parallel()

	db := &sqlx.DB{}

	if repo := NewBudgetRepository(db); repo == nil || repo.db != db {
		t.Fatalf("unexpected budget repo: %+v", repo)
	}
	if repo := NewCategoryRepository(db); repo == nil || repo.db != db {
		t.Fatalf("unexpected category repo: %+v", repo)
	}
	if repo := NewDeviceRepository(db); repo == nil || repo.db != db {
		t.Fatalf("unexpected device repo: %+v", repo)
	}
	if repo := NewFundRepository(db); repo == nil || repo.db != db {
		t.Fatalf("unexpected fund repo: %+v", repo)
	}
	if repo := NewNotificationRepository(db); repo == nil || repo.db != db {
		t.Fatalf("unexpected notification repo: %+v", repo)
	}
	if repo := NewReportRepository(db); repo == nil || repo.db != db {
		t.Fatalf("unexpected report repo: %+v", repo)
	}
	if repo := NewSavingsGoalRepository(db); repo == nil || repo.db != db {
		t.Fatalf("unexpected savings goal repo: %+v", repo)
	}
	if repo := NewSourcePaymentRepository(db); repo == nil || repo.db != db {
		t.Fatalf("unexpected source payment repo: %+v", repo)
	}
	if repo := NewTokenRepository(db); repo == nil || repo.db != db {
		t.Fatalf("unexpected token repo: %+v", repo)
	}
	if repo := NewTransactionRepository(db); repo == nil || repo.db != db {
		t.Fatalf("unexpected transaction repo: %+v", repo)
	}
	if repo := NewUserRepository(db); repo == nil || repo.db != db {
		t.Fatalf("unexpected user repo: %+v", repo)
	}
}
