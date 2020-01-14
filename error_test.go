package maperr

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorf(t *testing.T) {
	err := Errorf("some err with foo %d", 123)
	assert.EqualError(t, err, "some err with foo 123")
}

func TestNewError(t *testing.T) {
	err := NewError("some error")
	if !err.Equal(errors.New("some error")) {
		t.Fatal("expected errors to be equal")
	}
}

func TestCastError_FromErrorWithStatus(t *testing.T) {
	errWithStatus := WithStatus("BAD-REQUEST", http.StatusBadRequest)

	castedErr := castError(errWithStatus)

	if !errors.Is(castedErr, errWithStatus) {
		t.Fatalf("expected %s type %s got %s type %s", errWithStatus, reflect.TypeOf(errWithStatus), castedErr, reflect.TypeOf(castedErr))
	}
}

func TestFormattedError_Is(t *testing.T) {
	left := newFormattedError("foo failed: %s", "15644")
	right := newFormattedError("foo failed: %s", "8745616")

	if !left.Is(right) {
		t.Fatalf("expected %s to be the same error as %s", left, right)
	}
}

func TestFormattedError_Is_NotTheSameFormat(t *testing.T) {
	left := newFormattedError("foo failed: %s", "15644")
	right := newFormattedError("bar failed: %s", "8745616")

	if left.Is(right) {
		t.Fatalf("expected %s not be the same error as %s", left, right)
	}
}

func TestFormattedError_Equal(t *testing.T) {
	left := newFormattedError("foo failed: %s", "15644")
	right := newFormattedError("foo failed: %s", "8745616")

	if !left.Equal(right) {
		t.Fatalf("expected %s to be the same error as %s", left, right)
	}
}

func TestFormattedError_Equal_NotTheSameFormat(t *testing.T) {
	left := newFormattedError("foo failed: %s", "15644")
	right := newFormattedError("bar failed: %s", "8745616")

	if left.Equal(right) {
		t.Fatalf("expected %s not be the same error as %s", left, right)
	}
}

func TestFormattedError_Hashable(t *testing.T) {
	left := newFormattedError("foo failed: %s", "15644")

	if !errors.Is(left, left.Hashable()) {
		t.Fatalf("expected %s to be the same error as %s", left, left.Hashable())
	}
}

func TestFormattedError_Unwrap(t *testing.T) {
	err := newFormattedError("foo failed: %s", "15644")
	expected := "foo failed: 15644"

	if err.Unwrap().Error() != expected {
		t.Fatalf("expected %s to be the same error as %s", err.Unwrap().Error(), expected)
	}
}
