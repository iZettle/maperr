package maperr_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"

	"github.com/iZettle/maperr/v4"
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

func TestMap_ListMapper_FindAnyErrorInChain(t *testing.T) {
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
			NewListMapper().
			Append(errSecond, errors.New("this should be appended on Map()")),
	).Mapped(errChain, nil)

	if mappedErr == nil {
		t.Fatal("expected err got nil")
	}

	expected := "first error; second error; third error; forth error; fifth error; this should be appended on Map()"
	got := mappedErr.Error()
	if got != expected {
		t.Fatalf("expected %s got %s", expected, got)
	}
}
