package maperr

import (
	"errors"
	"fmt"
)

// Error which exposes a method that determines if the error
// can be considered equal or not
type Error interface {
	error
	Equal(error) bool
	Hashable() error
}

// Errorf returns an error which persists
func Errorf(format string, args ...interface{}) Error {
	return newFormattedError(format, args...)
}

// castError cast an error to maperr.Error when possible
// otherwise creates a new maperr.Error
func castError(err error) Error {
	if err == nil {
		return nil
	}

	var mapError Error
	if errors.As(err, &mapError) {
		return mapError
	}

	return NewError(err.Error())
}

// NewError instantiates an Error with no formatting
func NewError(errText string) Error {
	return Errorf(errText)
}

// formattedError is a error that holds the format
// from which the error was generated
type formattedError struct {
	format string
	args   []interface{}
	err    error
}

// newFormattedError return instance of formattedError
func newFormattedError(format string, args ...interface{}) formattedError {
	return formattedError{
		format: format,
		args:   args,
		err:    fmt.Errorf(format, args...),
	}
}

// Error return the actual error
func (fe formattedError) Error() string {
	return fe.err.Error()
}

// Unwrap return the actual error
func (fe formattedError) Unwrap() error {
	return fe.err
}

// Error return the hashable error
func (fe formattedError) Hashable() error {
	return fe.err
}

func (fe formattedError) Equal(err error) bool {
	if err == nil {
		return false
	}
	var ferr formattedError
	if errors.As(err, &ferr) {
		return fe.format == ferr.format
	}
	return fe.Error() == err.Error()
}
