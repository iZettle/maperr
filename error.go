package maperr

import "fmt"

// Errorf returns an error which persists
func Errorf(format string, args ...interface{}) error {
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

func (fe formattedError) EqualFormat(format string) bool {
	return fe.format == format
}
