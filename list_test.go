package maperr_test

import (
	"errors"
	"testing"

	"github.com/iZettle/maperr/v4"

	"github.com/stretchr/testify/assert"
)

func TestMap_MappedFormattedErrors(t *testing.T) {
	errTextLayerOneFailed := "foo %d"
	errLayerOneFailed := maperr.Errorf(errTextLayerOneFailed, 10)

	errTextLayerTwoFailed := "bar %d"
	errLayerTwoFailed := maperr.Errorf(errTextLayerTwoFailed, 20)

	type iteration struct {
		mapErr maperr.ListMapper
		err    error
	}
	tests := []struct {
		name        string
		iterations  []iteration
		expectedErr string
		defaultErr  error
	}{
		{
			name: "error going through three layers",
			iterations: []iteration{
				{
					mapErr: maperr.NewListMapper().
						Append(maperr.NewError("normal error"), errLayerOneFailed),
					err: errors.New("normal error"),
				},
				{
					mapErr: maperr.NewListMapper().
						Appendf(errTextLayerOneFailed, errLayerTwoFailed),
					err: errLayerOneFailed,
				},
				{
					mapErr: maperr.NewListMapper().
						Append(maperr.Errorf(errTextLayerTwoFailed), maperr.Errorf("abc")),
					err: errLayerTwoFailed,
				},
			},
			expectedErr: "normal error; foo 10; bar 20; abc",
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