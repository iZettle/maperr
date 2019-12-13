package maperr_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"

	"github.com/iZettle/maperr/v4"
)

func TestMultiErr_Mapped(t *testing.T) {
	first := errors.New("first")

	errListErrorOne := errors.New("errListErrorOne")
	errListErrorTwo := errors.New("errListErrorTwo")
	errListErrorThree := errors.New("errListErrorThree")

	errHashableErrorOne := errors.New("errHashableErrorOne")
	errHashableErrorTwo := errors.New("errHashableErrorTwo")

	errIgnore := errors.New("errIgnore")

	multipleMappers := maperr.NewMultiErr(
		maperr.NewListMapper().
			Append(errListErrorOne, errListErrorTwo).
			Append(errListErrorTwo, errListErrorThree),

		maperr.NewHashableMapper().
			Append(errHashableErrorOne, errHashableErrorTwo),

		maperr.NewIgnoreListMapper().
			Append(errIgnore),
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
			name:         "when errListErrorTwo is lastErr error, errListErrorThree is mapped and appended",
			mappedErrors: multipleMappers,
			givenError:   maperr.Combine(first, errListErrorOne, errListErrorTwo),
			expectedErr:  "first; errListErrorOne; errListErrorTwo; errListErrorThree",
		},
		{
			name:         "when errHashableErrorOne is lastErr error, errHashableErrorTwo is mapped and appended",
			mappedErrors: multipleMappers,
			givenError:   maperr.Combine(first, errHashableErrorOne),
			expectedErr:  "first; errHashableErrorOne; errHashableErrorTwo",
		},
		{
			name:         "errIgnore is ignored",
			mappedErrors: multipleMappers,
			givenError:   maperr.Combine(first, errIgnore),
			expectedErr:  "",
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

func TestMultiErr_LastMappedWithStatus(t *testing.T) {
	first := errors.New("first")
	second := errors.New("second")

	wrappedSecond := maperr.
		Append(first, second)

	mappedErrorsWithStatus := maperr.
		NewHashableMapper().
		Append(second, maperr.WithStatus("third", http.StatusInternalServerError))

	mappedErrorsWithoutStatus := maperr.
		NewHashableMapper().
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
			name:         "lastErr error was not found",
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
			name:         "lastErr error was mapped and wrapped",
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
			actualErr := maperr.
				NewMultiErr(test.mappedErrors).
				LastMappedWithStatus(test.givenError)
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

func TestMultiErr_MappedWithStatus(t *testing.T) {
	first := errors.New("first")
	second := errors.New("second")

	wrappedSecond := maperr.
		Append(first, second)

	mappedErrorsWithStatus := maperr.
		NewHashableMapper().
		Append(second, maperr.WithStatus("third", http.StatusInternalServerError))

	mappedErrorsWithoutStatus := maperr.
		NewHashableMapper().
		Append(second, errors.New("third"))

	type expected struct {
		status int
		err    string
	}
	tests := []struct {
		name         string
		mappedErrors maperr.Mapper
		givenError   error
		givenDefault error
		expected     expected
	}{
		{
			name:         "there was no error and nothing was provided",
			mappedErrors: mappedErrorsWithStatus,
			givenError:   nil,
			givenDefault: nil,
			expected: expected{
				err: "",
			},
		},
		{
			name:         "there was no error and a simple error was provided",
			mappedErrors: mappedErrorsWithStatus,
			givenError:   nil,
			givenDefault: errors.New("default error without a status code"),
			expected: expected{
				err: "",
			},
		},
		{
			name:         "there was no error and a error with status was provided as default",
			mappedErrors: mappedErrorsWithStatus,
			givenError:   nil,
			givenDefault: maperr.WithStatus("default error with a status code", http.StatusBadRequest),
			expected: expected{
				err: "",
			},
		},
		{
			name:         "lastErr error was not found and nothing was provided",
			mappedErrors: mappedErrorsWithStatus,
			givenError:   errors.New("not found"),
			givenDefault: nil,
			expected: expected{
				err: "",
			},
		},
		{
			name:         "lastErr error was not found and a simple error was provided",
			mappedErrors: mappedErrorsWithStatus,
			givenError:   errors.New("not found"),
			givenDefault: errors.New("default error without a status code"),
			expected: expected{
				status: http.StatusInternalServerError,
				err:    "default error without a status code",
			},
		},
		{
			name:         "lastErr error was not found and a error with status was provided as default",
			mappedErrors: mappedErrorsWithStatus,
			givenError:   errors.New("not found"),
			givenDefault: maperr.WithStatus("default error with a status code", http.StatusBadRequest),
			expected: expected{
				status: http.StatusBadRequest,
				err:    "default error with a status code",
			},
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
			name:         "lastErr error was mapped and wrapped",
			mappedErrors: mappedErrorsWithStatus,
			givenError:   wrappedSecond,
			expected: expected{
				status: http.StatusInternalServerError,
				err:    "third",
			},
		},
		{
			name:         "passed in error is nil, expect nil back",
			mappedErrors: maperr.NewIgnoreListMapper().Append(errors.New("test")),
			givenError:   nil,
			expected:     expected{},
		},
		{
			name:         "ignored with default",
			mappedErrors: maperr.NewIgnoreListMapper().Append(second),
			givenError:   wrappedSecond,
			givenDefault: errors.New("default error"),
			expected:     expected{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualErr := maperr.
				NewMultiErr(test.mappedErrors).
				MappedWithStatus(test.givenError, test.givenDefault)
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
		expected error
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
			expected: errors.New("three"),
		},
		{
			name:     "maperr.Error to be found is in the list",
			errList:  multierr.Combine(errors.New("one"), errors.New("two"), errors.New("three")),
			toFind:   errors.New("three"),
			expected: errors.New("three"),
		},
		{
			name: "formatted error to be found is in the list with a different id",
			errList: multierr.Combine(
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
