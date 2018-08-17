package maperr_test

import (
	"testing"

	"github.com/intelligentpos/maperr"
	"github.com/stretchr/testify/assert"
)

func TestErrorf(t *testing.T) {
	err := maperr.Errorf("some err with foo %d", 123)
	assert.EqualError(t, err, "some err with foo 123")
}
