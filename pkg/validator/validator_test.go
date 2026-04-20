package validator

import (
	"errors"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type samplePayload struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func TestValidateSuccess(t *testing.T) {
	t.Parallel()

	if err := Validate(samplePayload{Name: "Alice", Email: "alice@example.com"}); err != nil {
		t.Fatalf("Validate() unexpected error = %v", err)
	}
}

func TestValidateValidationErrorsUseJSONFieldNames(t *testing.T) {
	t.Parallel()

	err := Validate(samplePayload{})
	if err == nil {
		t.Fatal("expected validation error")
	}

	var fiberErr *fiber.Error
	if !errors.As(err, &fiberErr) {
		t.Fatalf("expected fiber error, got %T", err)
	}
	if fiberErr.Code != fiber.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", fiberErr.Code)
	}
	if fiberErr.Message != "Field 'name' failed on the 'required' tag; Field 'email' failed on the 'required' tag" {
		t.Fatalf("unexpected validation message: %q", fiberErr.Message)
	}
}

func TestValidateInvalidValidationInput(t *testing.T) {
	t.Parallel()

	err := Validate(nil)
	if err == nil {
		t.Fatal("expected validation setup error")
	}

	var fiberErr *fiber.Error
	if !errors.As(err, &fiberErr) {
		t.Fatalf("expected fiber error, got %T", err)
	}
	if fiberErr.Code != fiber.StatusInternalServerError {
		t.Fatalf("expected internal server error, got %d", fiberErr.Code)
	}
}
