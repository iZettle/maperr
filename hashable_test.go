package maperr_test

import (
	"errors"
	"testing"

	"github.com/iZettle/maperr/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

func TestMap_Mapped(t *testing.T) {
	layerOneFailed := errors.New("layer 1 failed")
	layerTwoFailed := errors.New("layer 2 failed")
	layerThreeFailed := errors.New("layer 3 failed")

	type iteration struct {
		mapErr maperr.HashableMapper
		err    error
	}
	tests := []struct {
		name        string
		iterations  []iteration
		expectedErr string
		defaultErr  error
	}{
		{
			name: "error is not mapped",
			iterations: []iteration{
				{
					mapErr: maperr.NewHashableMapper(),
					err:    errors.New("random error"),
				},
			},
			expectedErr: "random error",
		},
		{
			name: "error is not mapped",
			iterations: []iteration{
				{
					mapErr: maperr.NewHashableMapper(),
					err:    errors.New("random error"),
				},
			},
			expectedErr: "random error; some default description for generic error",
			defaultErr:  errors.New("some default description for generic error"),
		},
		{
			name: "error going through one layers",
			iterations: []iteration{
				{
					mapErr: maperr.NewHashableMapper().
						Append(layerOneFailed, layerTwoFailed),
					err: layerOneFailed,
				},
			},
			expectedErr: "layer 1 failed; layer 2 failed",
		},
		{
			name: "error going through three layers",
			iterations: []iteration{
				{
					mapErr: maperr.NewHashableMapper().
						Append(layerOneFailed, layerTwoFailed),
					err: layerOneFailed,
				},
				{
					mapErr: maperr.NewHashableMapper().
						Append(layerTwoFailed, layerThreeFailed),
					err: layerTwoFailed,
				},
			},
			expectedErr: "layer 1 failed; layer 2 failed; layer 3 failed",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualErr := test.iterations[0].err
			for _, iteration := range test.iterations {
				actualErr = maperr.NewMultiErr(iteration.mapErr).Mapped(actualErr, test.defaultErr)
			}
			if test.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.EqualError(t, actualErr, test.expectedErr)
			}
		})
	}
}

func TestMap_Mapped_FindAnyErrorInChain(t *testing.T) {
	errSecond := errors.New("second error")

	errChain := multierr.Combine(
		errors.New("first error"),
		errSecond,
		errors.New("third error"),
		errors.New("forth error"),
		errors.New("fifth error"),
	)

	mappedErr := maperr.NewMultiErr(
		maperr.
			NewHashableMapper().
			Append(errSecond, errors.New("this should be appended on mapErr()")),
	).Mapped(errChain, nil)

	if mappedErr == nil {
		t.Fatal("expected err got nil")
	}

	expected := "first error; second error; third error; forth error; fifth error; this should be appended on mapErr()"
	got := mappedErr.Error()
	if got != expected {
		t.Fatalf("expected %s got %s", expected, got)
	}
}
