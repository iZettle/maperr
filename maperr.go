package maperr

import (
	"errors"

	"go.uber.org/multierr"
)

// Mapper takes an error and return a matched error
type Mapper interface {
	Map(error) error
}

// HashableMapper simple implementation of Mapper which only works
// on hashable error keys
type HashableMapper map[error]error

// Map an error to another error
func (hm HashableMapper) Map(err error) error {
	err, ok := hm[err]
	if !ok {
		return nil
	}
	return err
}

// FormattedMapper maps formatted error strings to an error
type FormattedMapper map[string]error

// Map a formatted error to an error
func (hm FormattedMapper) Map(err error) error {
	fe, ok := err.(formattedError)
	if !ok {
		return nil
	}
	var matched error
	for errText, errMatchCandidate := range hm {
		if !fe.EqualFormat(errText) {
			continue
		}
		matched = errMatchCandidate
		break
	}
	return matched
}

type mapperList []Mapper

func (ml mapperList) Map(err error) error {
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
		mappers: mapperList(mapper),
	}
}

// Mapped appends the mapped error or a default one when is not found
func (m MultiErr) Mapped(err error, defaultErr error) error {
	if appendedErr := m.appendMapped(err); appendedErr != nil {
		return appendedErr
	}
	if err != nil && defaultErr != nil {
		return multierr.Append(err, defaultErr)
	}
	return err
}

// LastMapped return the last mapped error
func (m MultiErr) LastMapped(err error) error {
	var errorToMap = err
	previous := LastAppended(err)
	if previous != nil {
		errorToMap = previous
	}
	return m.mappers.Map(errorToMap)
}

// appendMapped append the mapped error to the given error
func (m MultiErr) appendMapped(err error) error {
	var errorToMap = err
	previous := LastAppended(err)
	if previous != nil {
		errorToMap = previous
	}
	mapped := m.mappers.Map(errorToMap)
	if mapped == nil {
		return nil
	}
	return multierr.Append(err, mapped)
}

// ErrorWithStatusProvider defines an error which also has an http status defined
type ErrorWithStatusProvider interface {
	error
	Status() int
}

// LastMappedWithStatus return the last mapped error with the associated http status
func (m MultiErr) LastMappedWithStatus(err error) ErrorWithStatusProvider {
	lastMapped := m.LastMapped(err)
	if lastMapped == nil {
		return nil
	}
	statusErr, ok := lastMapped.(ErrorWithStatusProvider)
	if !ok {
		return nil
	}
	return statusErr
}

type errorWithStatus struct {
	err    error
	status int
}

func (ews errorWithStatus) Status() int {
	return ews.status
}

func (ews errorWithStatus) Error() string {
	return ews.err.Error()
}

// WithStatus return an error with an associated status
func WithStatus(err string, status int) error {
	return &errorWithStatus{
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
