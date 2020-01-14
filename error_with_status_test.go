package maperr

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrorWithStatus_Hashable(t *testing.T) {
	errMissingField := errors.New("MISSING_FIELD")
	errWithStatus := newErrorWithStatus(errMissingField, errors.New("could not fill struct"), http.StatusBadRequest)

	if !errors.Is(errWithStatus, errWithStatus.Hashable()) {
		t.Fatalf("expected %s to be the same error as %s", errWithStatus, errWithStatus.Hashable())
	}
}

func TestErrorWithStatus_Is(t *testing.T) {
	errMissingField := errors.New("MISSING_FIELD")
	left := newErrorWithStatus(errMissingField, errors.New("could not fill struct"), http.StatusBadRequest)
	right := newErrorWithStatus(errMissingField, errors.New("could not fill struct"), http.StatusBadRequest)

	if !left.Is(right) {
		t.Fatalf("expected %s to be the same error as %s", left, right)
	}
}

func TestErrorWithStatus_Is_NotTheSameError(t *testing.T) {
	left := newErrorWithStatus(errors.New("MISSING_FIELD_ONE"), errors.New("could not fill struct"), http.StatusBadRequest)
	right := newErrorWithStatus(errors.New("MISSING_FIELD_TWO"), errors.New("could not fill struct"), http.StatusBadRequest)

	if left.Is(right) {
		t.Fatalf("expected %s to be the different error than %s", left, right)
	}
}

func TestErrorWithStatus_Is_NilError(t *testing.T) {
	errWithStatus := newErrorWithStatus(nil, errors.New("could not fill struct"), http.StatusBadRequest)

	if errWithStatus.Is(nil) {
		t.Fatalf("expected %s to be the different error than nil", errWithStatus)
	}
}

func TestErrorWithStatus_Equal(t *testing.T) {
	errMissingField := errors.New("MISSING_FIELD")
	left := newErrorWithStatus(errMissingField, errors.New("could not fill struct"), http.StatusBadRequest)
	right := newErrorWithStatus(errMissingField, errors.New("could not fill struct"), http.StatusBadRequest)

	if !left.Equal(right) {
		t.Fatalf("expected %s to be the same error as %s", left, right)
	}
}

func TestErrorWithStatus_Equal_NotTheSameError(t *testing.T) {
	left := newErrorWithStatus(errors.New("MISSING_FIELD_ONE"), errors.New("could not fill struct"), http.StatusBadRequest)
	right := newErrorWithStatus(errors.New("MISSING_FIELD_TWO"), errors.New("could not fill struct"), http.StatusBadRequest)

	if left.Equal(right) {
		t.Fatalf("expected %s to be the different error than %s", left, right)
	}
}

func TestErrorWithStatus_Equal_NilError(t *testing.T) {
	errWithStatus := newErrorWithStatus(nil, errors.New("could not fill struct"), http.StatusBadRequest)

	if errWithStatus.Equal(nil) {
		t.Fatalf("expected %s to be the different error than nil", errWithStatus)
	}
}
