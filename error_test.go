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
