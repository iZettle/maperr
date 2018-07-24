package maperr_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/intelligentpos/maperr"
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
					mapErr: maperr.HashableMapper{},
					err:    errors.New("random error"),
				},
			},
			expectedErr: "random error",
		},
		{
			name: "error is not mapped",
			iterations: []iteration{
				{
					mapErr: maperr.HashableMapper{},
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
					mapErr: maperr.HashableMapper{
						layerOneFailed: layerTwoFailed,
					},
					err: layerOneFailed,
				},
			},
			expectedErr: "layer 1 failed; layer 2 failed",
		},
		{
			name: "error going through three layers",
			iterations: []iteration{
				{
					mapErr: maperr.HashableMapper{
						layerOneFailed: layerTwoFailed,
					},
					err: layerOneFailed,
				},
				{
					mapErr: maperr.HashableMapper{
						layerTwoFailed: layerThreeFailed,
					},
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

func TestMap_LastAppended(t *testing.T) {
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
			name:             "several errors",
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

func TestMap_LastMapped(t *testing.T) {
	first := errors.New("first")
	second := errors.New("second")
	third := errors.New("third")
	wrappedSecond := multierr.Append(first, second)
	mappedErrors := maperr.HashableMapper{
		second: third,
	}
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

func TestMap_LastMappedWithStatus(t *testing.T) {
	first := errors.New("first")
	second := errors.New("second")
	wrappedSecond := multierr.Append(first, second)
	mappedErrorsWithStatus := maperr.HashableMapper{
		second: maperr.WithStatus("third", http.StatusInternalServerError),
	}
	mappedErrorsWithoutStatus := maperr.HashableMapper{
		second: errors.New("third"),
	}
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
			name:         "mapped errors without an http status",
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
			errList:  multierr.Append(errors.New("one"), errors.New("two")),
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
