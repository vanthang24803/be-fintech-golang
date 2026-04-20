package snowflake

import (
	"testing"
)

func TestGenerateID_NonZero(t *testing.T) {
	id := GenerateID()
	if id == 0 {
		t.Fatal("expected non-zero ID")
	}
}

func TestGenerateID_Unique(t *testing.T) {
	seen := make(map[int64]bool)
	for i := 0; i < 100; i++ {
		id := GenerateID()
		if seen[id] {
			t.Fatalf("duplicate ID generated: %d", id)
		}
		seen[id] = true
	}
}

func TestInit(t *testing.T) {
	Init(1)
	id := GenerateID()
	if id == 0 {
		t.Fatal("expected non-zero ID after Init")
	}
}
