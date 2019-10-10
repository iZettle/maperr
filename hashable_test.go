package maperr_test

import (
	"errors"
	"testing"

	"github.com/iZettle/maperr/v4"
	"github.com/stretchr/testify/assert"
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
