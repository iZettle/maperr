package maperr_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iZettle/maperr/v4"
)

func TestErrorf(t *testing.T) {
	err := maperr.Errorf("some err with foo %d", 123)
	assert.EqualError(t, err, "some err with foo 123")
}

func TestNewError(t *testing.T) {
	err := maperr.NewError("some error")
	if !err.Equal(errors.New("some error")) {
		t.Fatal("expected errors to be equal")
	}
}
