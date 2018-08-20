package maperr

import "fmt"

// Error which exposes a method that determines if the error
// can be considered equal or not
type Error interface {
	error
	Equal(Error) bool
}

// Errorf returns an error which persists
func Errorf(format string, args ...interface{}) Error {
	return newFormattedError(format, args...)
}

// formattedError is a error that holds the format
// from which the error was generated
type formattedError struct {
	format string
	args   []interface{}
	error  error
}

// newFormattedError return instance of formattedError
func newFormattedError(format string, args ...interface{}) formattedError {
	return formattedError{
		format: format,
		args:   args,
		error:  fmt.Errorf(format, args...),
	}
}

// Error return the actual error
func (fe formattedError) Error() string {
	return fe.error.Error()
}

func (fe formattedError) Equal(err Error) bool {
	formattedErr, ok := err.(formattedError)
	if !ok {
		return false
	}
	return fe.format == formattedErr.format
}
