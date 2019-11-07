package maperr_test

import (
	"errors"
	"testing"

	"go.uber.org/multierr"

	"github.com/iZettle/maperr/v4"
)

func Test_IgnoreListMapper_Map_IgnoreErrorFound(t *testing.T) {
	errLayerOneFailed := errors.New("maperr error")
	stdLibraryError := errors.New("maperr error")
	errTextLayerTwoFailed := "bar %d"

	tests := []struct {
		name  string
		given error
	}{
		{
			name:  "error ignored",
			given: errLayerOneFailed,
		},
		{
			name:  "std library error ignored",
			given: stdLibraryError,
		},
		{
			name:  "formatted error ignored",
			given: maperr.Errorf(errTextLayerTwoFailed, "foo"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// append the error at the end of a series of errors to ensure
			// that we are comparing against the last error
			chainOfErrors := maperr.Combine(
				errors.New("first"),
				errors.New("second"),
				errors.New("third"),
				test.given,
			)

			ignoreMapper := maperr.NewIgnoreListMapper().
				Append(errLayerOneFailed).
				Appendf(errTextLayerTwoFailed)

			res := ignoreMapper.Map(chainOfErrors)
			if res == nil {
				t.Fatal("expected a map result object")
			}
			if res.Previous().Error() != chainOfErrors.Error() {
				t.Fatalf("expected %s got %s", chainOfErrors.Error(), res.Previous().Error())
			}
			if res.Last() != nil {
				t.Fatalf("expected nil got %s", res.Last().Error())
			}
			if res.Apply() != nil {
				t.Fatalf("expected nil got %s", res.Apply().Error())
			}
		})
	}
}

func Test_IgnoreListMapper_Map_IgnoreErrorNotFound(t *testing.T) {
	errLayerOneFailed := errors.New("layer 1 failed")
	errTextLayerTwoFailed := "bar %d"

	tests := []struct {
		name  string
		given error
	}{
		{
			name:  "not found",
			given: errors.New("not found"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// append the error at the end of a series of errors to ensure
			// that we are comparing against the last error
			chainOfErrors := maperr.Combine(
				errors.New("first"),
				errors.New("second"),
				errors.New("third"),
				test.given,
			)

			ignoreMapper := maperr.NewIgnoreListMapper().
				Append(errLayerOneFailed).
				Appendf(errTextLayerTwoFailed)

			res := ignoreMapper.Map(chainOfErrors)
			if res != nil {
				t.Fatal("expected nil got a map result object")
			}
		})
	}
}

func Test_IgnoreListMapper_Map_WithMultiErr(t *testing.T) {
	errLayerOneFailed := errors.New("layer 1 failed")
	errTextLayerTwoFailed := "bar %d"

	tests := []struct {
		name  string
		given error
	}{
		{
			name:  "error ignored",
			given: errLayerOneFailed,
		},
		{
			name:  "formatted error ignored",
			given: maperr.Errorf(errTextLayerTwoFailed, "foo"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// append the error at the end of a series of errors to ensure
			// that we are comparing against the last error
			chainOfErrors := maperr.Combine(
				errors.New("first"),
				errors.New("second"),
				errors.New("third"),
				test.given,
			)

			ignoreMapper := maperr.NewIgnoreListMapper().
				Append(errLayerOneFailed).
				Appendf(errTextLayerTwoFailed)

			mappedErr := maperr.NewMultiErr(ignoreMapper).
				Mapped(chainOfErrors, errors.New("default error"))
			if mappedErr != nil {
				t.Fatalf("expected nil got err %s", mappedErr.Error())
			}
		})
	}
}

func TestMap_IgnoreListMapper_FindAnyErrorInChain(t *testing.T) {
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
			NewIgnoreListMapper().
			Append(errSecond),
	).Mapped(errChain, nil)

	if mappedErr != nil {
		t.Fatal("expected nil got err")
	}
}
