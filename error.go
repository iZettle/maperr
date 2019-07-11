package maperr

import "fmt"

// Error which exposes a method that determines if the error
// can be considered equal or not
type Error interface {
	error
	Equal(Error) bool
	Hashable() error
}

// Errorf returns an error which persists
func Errorf(format string, args ...interface{}) Error {
	return newFormattedError(format, args...)
}

// CastError cast an error to maperr.Error when possible
// otherwise creates a new maperr.Error
func CastError(err error) Error {
	if err == nil {
		return nil
	}
	if mapError, ok := err.(Error); ok {
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

// Error return the hashable error
func (fe formattedError) Hashable() error {
	return fe.err
}

func (fe formattedError) Equal(err Error) bool {
	formattedErr, ok := err.(formattedError)
	if !ok {
		return false
	}
	return fe.format == formattedErr.format
}
