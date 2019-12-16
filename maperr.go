package maperr

import (
	"errors"
	"net/http"

	"go.uber.org/multierr"
)

// MapResult is an interface that defines the result of a Map operation
type MapResult interface {
	Previous() error
	Last() error
	Apply() error
}

// Mapper takes an error and return a MapResult
type Mapper interface {
	Map(error) MapResult
}

type mapperList []Mapper

func (ml mapperList) Map(err error) MapResult {
	for k := range ml {
		if mapped := ml[k].Map(err); mapped != nil {
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
	if res := m.mappers.Map(err); res != nil {
		return res.Apply()
	}
	if defaultErr != nil {
		return Append(err, defaultErr)
	}
	return err
}

// lastMapped return the last mapped error
func (m MultiErr) lastMapped(err error) MapResult {
	res := m.mappers.Map(err)
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

// MappedWithStatus return the last mapped error with the associated http status
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

	// if the last appended error was ignored we have to map it to nil
	lastMappedResult := m.lastMapped(err)
	if _, ok := lastMappedResult.(ignoreStrategy); ok {
		return nil
	}

	if lastMappedResult == nil && err != nil {
		var defaultErrWithStatus ErrorWithStatusProvider
		if errors.As(defaultErr, &defaultErrWithStatus) {
			return newErrorWithStatus(defaultErrWithStatus, err, defaultErrWithStatus.Status())
		}
		if defaultErr != nil {
			return newErrorWithStatus(defaultErr, err, http.StatusInternalServerError)
		}
	}
	if lastMappedResult == nil {
		return nil
	}

	lastMapped := lastMappedResult.Last()
	var statusErr ErrorWithStatusProvider
	if errors.As(lastMapped, &statusErr) {
		return statusErr
	}

	return nil
}

// LastMappedWithStatus return the last mapped error with the associated http status
// Deprecated: consider using MappedWithStatus() instead, as encourages to specify a default error
func (m MultiErr) LastMappedWithStatus(err error) ErrorWithStatusProvider {
	return m.MappedWithStatus(err, nil)
}

type errorWithStatus struct {
	err    error
	status int
	cause  error
}

func newErrorWithStatus(err, cause error, status int) errorWithStatus {
	return errorWithStatus{
		err:    err,
		status: status,
		cause:  cause,
	}
}

func (ews errorWithStatus) Status() int {
	return ews.status
}

func (ews errorWithStatus) Unwrap() error {
	return ews.cause
}

func (ews errorWithStatus) Error() string {
	return ews.err.Error()
}

// WithStatus return an error with an associated status
func WithStatus(err string, status int) error {
	return errorWithStatus{
		err:    errors.New(err),
		status: status,
	}
}

// LastAppended return the last error appended as multierr
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
