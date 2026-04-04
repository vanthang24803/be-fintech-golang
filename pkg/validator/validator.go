package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register tag name to be the same as json tag
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate performs validation on a struct
func Validate(i interface{}) error {
	err := validate.Struct(i)
	if err == nil {
		return nil
	}

	var invalidErr *validator.InvalidValidationError
	if errors.As(err, &invalidErr) {
		return fiber.NewError(fiber.StatusInternalServerError, "Validation setup error")
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errMsgs := make([]string, 0, len(validationErrors))
		for _, err := range validationErrors {
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' failed on the '%s' tag", err.Field(), err.Tag()))
		}
		return fiber.NewError(fiber.StatusBadRequest, strings.Join(errMsgs, "; "))
	}

	return fiber.NewError(fiber.StatusBadRequest, err.Error())
}
