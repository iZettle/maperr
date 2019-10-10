package maperr_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/iZettle/maperr/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

func TestMultiErr_Mapped(t *testing.T) {
	first := maperr.NewError("first")
	second := maperr.NewError("second")
	third := maperr.NewError("third")
	forth := maperr.NewError("forth")
	fifth := maperr.NewError("fifth")

	multipleMappers := maperr.NewMultiErr(
		maperr.NewIgnoreListMapper().
			Append(second),

		maperr.NewListMapper().
			Append(second, third).
			Append(third, forth),

		maperr.NewHashableMapper().
			Append(forth, fifth),
	)

	tests := []struct {
		name         string
		mappedErrors maperr.MultiErr
		givenError   error
		expectedErr  string
	}{
		{
			name:         "error is nil",
			mappedErrors: multipleMappers,
			givenError:   nil,
			expectedErr:  "",
		},
		{
			name:         "second is ignored",
			mappedErrors: multipleMappers,
			givenError:   maperr.Append(first, second),
			expectedErr:  "",
		},
		{
			name:         "when third is last error, forth is mapped and appended",
			mappedErrors: multipleMappers,
			givenError:   maperr.Combine(first, second, third),
			expectedErr:  "first; second; third; forth",
		},
		{
			name:         "when forth is last error, fifth is mapped and appended",
			mappedErrors: multipleMappers,
			givenError:   maperr.Combine(first, second, third, forth),
			expectedErr:  "first; second; third; forth; fifth",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualErr := test.mappedErrors.Mapped(test.givenError, errors.New("default"))
			if test.expectedErr != "" {
				assert.EqualError(t, actualErr, test.expectedErr)
			} else {
				assert.NoError(t, actualErr)
			}
		})
	}
}

func TestMultiErr_LastMapped(t *testing.T) {
	first := errors.New("first")
	second := errors.New("second")
	third := errors.New("third")
	wrappedSecond := maperr.Append(first, second)
	mappedErrors := maperr.NewHashableMapper().
		Append(second, third)
	tests := []struct {
		name         string
		mappedErrors maperr.HashableMapper
		givenError   error
		expectedErr  string
	}{
		{
			name:         "last error was not found",
			mappedErrors: mappedErrors,
			givenError:   errors.New("not found"),
		},
		{
			name:         "last error was mapped and wrapped",
			mappedErrors: mappedErrors,
			givenError:   wrappedSecond,
			expectedErr:  "third",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualErr := maperr.NewMultiErr(test.mappedErrors).LastMapped(test.givenError)
			if test.expectedErr != "" {
				assert.EqualError(t, actualErr, test.expectedErr)
			} else {
				assert.NoError(t, actualErr)
			}
		})
	}
}

func TestMultiErr_LastMappedWithStatus(t *testing.T) {
	first := errors.New("first")
	second := errors.New("second")
	wrappedSecond := maperr.
		Append(first, second)
	mappedErrorsWithStatus := maperr.NewHashableMapper().
		Append(second, maperr.WithStatus("third", http.StatusInternalServerError))
	mappedErrorsWithoutStatus := maperr.NewHashableMapper().
		Append(second, errors.New("third"))
	type expected struct {
		status int
		err    string
	}
	tests := []struct {
		name         string
		mappedErrors maperr.HashableMapper
		givenError   error
		expected     expected
	}{
		{
			name:         "last error was not found",
			mappedErrors: mappedErrorsWithStatus,
			givenError:   errors.New("not found"),
		},
		{
			name:         "mapped errorPairs without an http status",
			mappedErrors: mappedErrorsWithoutStatus,
			givenError:   wrappedSecond,
			expected: expected{
				status: 0,
				err:    "",
			},
		},
		{
			name:         "last error was mapped and wrapped",
			mappedErrors: mappedErrorsWithStatus,
			givenError:   wrappedSecond,
			expected: expected{
				status: http.StatusInternalServerError,
				err:    "third",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualErr := maperr.NewMultiErr(test.mappedErrors).LastMappedWithStatus(test.givenError)
			if test.expected.err != "" {
				assert.EqualError(t, actualErr, test.expected.err)
			} else {
				assert.NoError(t, actualErr)
			}

			if actualErr != nil {
				assert.Equal(t, actualErr.Status(), test.expected.status)
			}
		})
	}
}

func TestHasError(t *testing.T) {
	tests := []struct {
		name     string
		errList  error
		toFind   string
		expected bool
	}{
		{
			name:     "err is not a list",
			errList:  errors.New("foo"),
			expected: false,
		},
		{
			name:     "err is not found",
			errList:  maperr.Append(errors.New("one"), errors.New("two")),
			toFind:   "three",
			expected: false,
		},
		{
			name:     "err is not found",
			errList:  multierr.Combine(errors.New("one"), errors.New("two"), errors.New("three")),
			toFind:   "three",
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, maperr.HasError(test.errList, test.toFind))
		})
	}
}

func TestLastAppended(t *testing.T) {
	layerOneFailed := errors.New("layer 1 failed")
	layerTwoFailed := errors.New("layer 2 failed")
	layerThreeFailed := errors.New("layer 3 failed")
	tests := []struct {
		name             string
		err              error
		expectedPrevious error
	}{
		{
			name:             "empty error",
			err:              nil,
			expectedPrevious: nil,
		},
		{
			name:             "one error",
			err:              layerOneFailed,
			expectedPrevious: layerOneFailed,
		},
		{
			name:             "several errorPairs",
			err:              multierr.Combine(layerOneFailed, layerTwoFailed, layerThreeFailed),
			expectedPrevious: layerThreeFailed,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := maperr.LastAppended(test.err)
			assert.Equal(t, test.expectedPrevious, actual)
		})
	}
}

func Test_HasEqual(t *testing.T) {
	tests := []struct {
		name     string
		errList  error
		toFind   error
		expected maperr.Error
	}{
		{
			name:     "error to be found is nil",
			errList:  errors.New("foo"),
			toFind:   nil,
			expected: nil,
		},
		{
			name:     "list is nil",
			errList:  nil,
			toFind:   errors.New("three"),
			expected: nil,
		},
		{
			name:     "error to be found is not in list",
			errList:  maperr.Append(errors.New("one"), errors.New("two")),
			toFind:   errors.New("three"),
			expected: nil,
		},
		{
			name:     "error to be found is in the list",
			errList:  multierr.Combine(errors.New("one"), errors.New("two"), errors.New("three")),
			toFind:   errors.New("three"),
			expected: maperr.NewError("three"),
		},
		{
			name:     "maperr.Error to be found is in the list",
			errList:  multierr.Combine(errors.New("one"), errors.New("two"), errors.New("three")),
			toFind:   maperr.NewError("three"),
			expected: maperr.NewError("three"),
		},
		{
			name:     "formatted error to be found is in the list with a different id",
			errList:  multierr.Combine(
				errors.New("one"),
				maperr.Errorf("this is a formatted error %d", 12345),
				errors.New("two"),
				errors.New("three")),
			toFind:   maperr.Errorf("this is a formatted error %d", 98765),
			expected: maperr.Errorf("this is a formatted error %d", 12345),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := maperr.HasEqual(test.errList, test.toFind)

			if got == nil && test.expected != nil {
				t.Fatalf("expected %s got nil", test.expected.Error())
			}

			if got != nil && test.expected == nil {
				t.Fatalf("expected nil got %s", got.Error())
			}

			if got != nil && test.expected != nil && got.Error() != test.expected.Error() {
				t.Fatalf("expected %s got %s", test.expected.Error(), got.Error())
			}
		})
	}
}
