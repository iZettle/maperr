package maperr

import (
	"errors"
)

// WithStatus return an error with an associated status
func WithStatus(err string, status int) error {
	return errorWithStatus{
		err:    errors.New(err),
		status: status,
	}
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

func (ews errorWithStatus) Hashable() error {
	return ews
}

// Is is an alias for Equal added to support go 1.13 errors
func (ews errorWithStatus) Is(err error) bool {
	return ews.Equal(err)
}

func (ews errorWithStatus) Equal(err error) bool {
	if err == nil {
		return false
	}
	var errWithStatus errorWithStatus
	if errors.As(err, &errWithStatus) {
		return errors.Is(ews.err, errWithStatus.err)
	}
	return false
}
