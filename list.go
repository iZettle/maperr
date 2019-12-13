package maperr

import (
	"go.uber.org/multierr"
)

// PairErrors holds a pair of errorPairs
type PairErrors struct {
	err   Error
	match Error
}

// ListMapper maps not hashable or formatted errorPairs
type ListMapper struct {
	errorPairs []PairErrors
}

// NewListMapper return a new ListMapper
func NewListMapper() ListMapper {
	return ListMapper{}
}

// Appendf append a format to error association
func (lm ListMapper) Appendf(format string, match error) ListMapper {
	return lm.Append(Errorf(format), castError(match))
}

// Append append an error to error association
func (lm ListMapper) Append(err, match error) ListMapper {
	lm.errorPairs = append(lm.errorPairs,
		PairErrors{
			err:   castError(err),
			match: castError(match),
		})
	return lm
}

// mapErr a formatted error to an error
func (lm ListMapper) mapErr(err error) mapResult {
	errorsToMap := []error{
		err,
	}
	if errList := multierr.Errors(err); len(errList) > 0 {
		errorsToMap = errList
	}

	for i := len(errorsToMap) - 1; i >= 0; i-- {
		comparableErr := castError(errorsToMap[i])
		for k := range lm.errorPairs {
			if comparableErr.Equal(lm.errorPairs[k].err) {
				return newAppendStrategy(err, lm.errorPairs[k].match)
			}
		}
	}

	return nil
}

// appendStrategy holds data for an ignore strategy
type appendStrategy struct {
	previousErr error
	lastErr     error
}

// newAppendStrategy instantiates a new appendStrategy
func newAppendStrategy(previous, last error) appendStrategy {
	return appendStrategy{previousErr: previous, lastErr: last}
}

// previous returns the error that we want to append to
func (as appendStrategy) previous() error {
	return as.previousErr
}

// last returns the error that we are appending
func (as appendStrategy) last() error {
	return as.lastErr
}

// apply the append strategy by appending previousErr to lastErr
func (as appendStrategy) apply() error {
	if as.lastErr == nil {
		return nil
	}
	return Append(as.previousErr, as.lastErr)
}
