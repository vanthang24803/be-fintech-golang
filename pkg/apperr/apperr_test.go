package apperr

import (
	"errors"
	"fmt"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		err     error
		wantMsg string
	}{
		{"NotFound", ErrNotFound, "not found"},
		{"InsufficientBalance", ErrInsufficientBalance, "insufficient balance"},
		{"Conflict", ErrConflict, "resource already exists"},
		{"Unauthorized", ErrUnauthorized, "unauthorized"},
		{"InvalidInput", ErrInvalidInput, "invalid input"},
		{"Internal", ErrInternal, "internal server error"},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.wantMsg {
				t.Fatalf("expected %q, got %q", tt.wantMsg, tt.err.Error())
			}

			wrapped := fmt.Errorf("wrapped: %w", tt.err)
			if !errors.Is(wrapped, tt.err) {
				t.Fatalf("errors.Is should match wrapped error for %s", tt.name)
			}
		})
	}
}

func TestSentinelErrors_Distinct(t *testing.T) {
	t.Parallel()

	all := []error{ErrNotFound, ErrInsufficientBalance, ErrConflict, ErrUnauthorized, ErrInvalidInput, ErrInternal}
	for i, a := range all {
		for j, b := range all {
			if i != j && errors.Is(a, b) {
				t.Fatalf("errors[%d] and errors[%d] should be distinct but errors.Is returned true", i, j)
			}
		}
	}
}
