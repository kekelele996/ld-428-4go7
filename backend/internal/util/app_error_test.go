package util

import (
	"errors"
	"testing"
)

func TestAppErrorUnwrap(t *testing.T) {
	base := errors.New("db down")
	err := Wrap(50000, "服务器内部错误", base)
	if !errors.Is(err, base) {
		t.Fatal("Wrap should preserve underlying error")
	}
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Code != 50000 {
		t.Fatalf("unexpected app error: %+v", appErr)
	}
}

func TestIsAppError(t *testing.T) {
	cases := []struct {
		err  error
		code int
		want bool
	}{
		{NewAppError(40400, "not found", nil), 40400, true},
		{Wrap(40900, "conflict", errors.New("x")), 40900, true},
		{errors.New("x"), 40900, false},
	}
	for _, c := range cases {
		if got := IsAppError(c.err, c.code); got != c.want {
			t.Fatalf("IsAppError(%v,%d)=%v want %v", c.err, c.code, got, c.want)
		}
	}
}
