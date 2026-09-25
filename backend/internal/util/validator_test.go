package util

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestTranslateError(t *testing.T) {
	type req struct {
		Name string `validate:"required"`
		Age  int    `validate:"min=18"`
	}
	validate := validator.New()
	err := validate.Struct(req{})
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		t.Fatal("expected validation errors")
	}
	msg := TranslateError(err)
	if msg == "" || msg == err.Error() {
		t.Fatalf("expected translated message, got %q", msg)
	}
}
