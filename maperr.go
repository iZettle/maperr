package maperr

import (
	"errors"
	"net/http"

	"go.uber.org/multierr"
)

// mapResult is an interface that defines the result of a mapErr operation
type mapResult interface {
	previous() error
	last() error
	apply() error
}

// Mapper takes an error and return a mapResult
type Mapper interface {
	mapErr(error) mapResult
}

type mapperList []Mapper

func (ml mapperList) mapErr(err error) mapResult {
	for k := range ml {
		if mapped := ml[k].mapErr(err); mapped != nil {
			return mapped
		}
	}
	return nil
}

// MultiErr an error to another error
type MultiErr struct {
	mappers mapperList
}

// NewMultiErr return a new instance of MultiErr
func NewMultiErr(mapper ...Mapper) MultiErr {
	return MultiErr{
		mappers: mapper,
	}
}

// Mapped appends the mapped error or a default one when is not found
func (m MultiErr) Mapped(err, defaultErr error) error {
	if err == nil {
		return nil
	}
	if res := m.mappers.mapErr(err); res != nil {
		return res.apply()
	}
	if defaultErr != nil {
		return Append(err, defaultErr)
	}
	return err
}

// lastMapped return the lastErr mapped error
func (m MultiErr) lastMapped(err error) mapResult {
	res := m.mappers.mapErr(err)
	if res == nil {
		return nil
	}
	return res
}

// Default error with statuses
var (
	WithStatusBadRequest          = WithStatus(http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	WithStatusInternalServerError = WithStatus(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
)

// ErrorWithStatusProvider defines an error which also has an http status defined
type ErrorWithStatusProvider interface {
	error
	Status() int
	Unwrap() error
}

// MappedWithStatus return the lastErr mapped error with the associated http status
// You can optionally provide a default error in case that will be returned if the error has not been mapped
//
// defaultErr == nil                      returns the ErrorWithStatusProvider only if error is mapped
//                                        alias for LastMappedWithStatus(err)
//
// defaultErr.(ErrorWithStatusProvider)   allow you to specify a status code
//                                        e.g.: maperr.WithStatus("USER_ERROR", http.StatusBadRequest)
//
// defaultErr.(error)                     will cast to a ErrorWithStatusProvider with http.StatusInternalServerError
func (m MultiErr) MappedWithStatus(err, defaultErr error) ErrorWithStatusProvider {
	if err == nil {
		return nil
	}

	lastMappedResult := m.lastMapped(err)

	// when the mapped error comes from the "ignore list" we can exit early
	if _, ok := lastMappedResult.(ignoreStrategy); ok {
		return nil
	}

	// when have an error that could not be mapped, we use the defaultErr parameter instead
	if lastMappedResult == nil && err != nil {
		if defaultStatusErr := appendCauseToErrWithStatus(defaultErr, err); defaultStatusErr != nil {
			return defaultStatusErr
		}
		if defaultErr != nil {
			return newErrorWithStatus(defaultErr, err, http.StatusInternalServerError)
		}
	}
	if lastMappedResult == nil {
		return nil
	}

	lastMapped := lastMappedResult.last()
	if statusErr := appendCauseToErrWithStatus(lastMapped, err); statusErr != nil {
		return statusErr
	}

	return nil
}

func appendCauseToErrWithStatus(err, cause error) ErrorWithStatusProvider {
	var errWithStatus errorWithStatus
	if !errors.As(err, &errWithStatus) {
		return nil
	}
	errWithStatus.cause = cause

	return errWithStatus
}

// LastMappedWithStatus return the lastErr mapped error with the associated http status
// Deprecated: consider using MappedWithStatus() instead, as encourages to specify a default error
func (m MultiErr) LastMappedWithStatus(err error) ErrorWithStatusProvider {
	return m.MappedWithStatus(err, nil)
}

// LastAppended return the lastErr error appended as multierr
func LastAppended(err error) error {
	if errList := multierr.Errors(err); len(errList) > 0 {
		return errList[len(errList)-1]
	}
	return nil
}

// HasEqual find the first equal error on a chain of errors
// and returns it
func HasEqual(chain, err error) Error {
	mapError := castError(err)

	multiErrList := multierr.Errors(chain)
	if len(multiErrList) == 0 {
		return nil
	}

	for _, wrapped := range multiErrList {
		wrappedMapError := castError(wrapped)
		if wrappedMapError.Equal(mapError) {
			return wrappedMapError
		}
	}
	return nil
}

// HasError checks if an error has been wrapped
func HasError(err error, errText string) bool {
	multiErrList := multierr.Errors(err)
	if len(multiErrList) == 0 {
		return false
	}
	found := false
	for _, wrapped := range multiErrList {
		if wrapped.Error() == errText {
			found = true
			break
		}
	}
	return found
}

// Append appends the given errors together. Either value may be nil.
func Append(left, right error) error {
	return multierr.Append(left, right)
}

// Combine combines the passed errors into a single error.
func Combine(errList ...error) error {
	return multierr.Combine(errList...)
}
