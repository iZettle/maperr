package maperr

import (
	"go.uber.org/multierr"
)

// IgnoreListMapper is a Mapper that allow to specify a list of error that we
// want to ignore
type IgnoreListMapper struct {
	list []Error
}

// NewIgnoreListMapper create an instance of IgnoreListMapper
func NewIgnoreListMapper() IgnoreListMapper {
	return IgnoreListMapper{}
}

// Appendf appends a formatted error that we want to ignore
func (lm IgnoreListMapper) Appendf(format string) IgnoreListMapper {
	return lm.Append(Errorf(format))
}

// Append appends an error that we want to ignore
func (lm IgnoreListMapper) Append(err error) IgnoreListMapper {
	lm.list = append(lm.list, castError(err))
	return lm
}

// mapErr an error to an ignore strategy
func (lm IgnoreListMapper) mapErr(err error) mapResult {
	errorsToMap := []error{
		err,
	}
	if errList := multierr.Errors(err); len(errList) > 0 {
		errorsToMap = errList
	}

	for i := len(errorsToMap) - 1; i >= 0; i-- {
		comparableErr := castError(errorsToMap[i])
		for k := range lm.list {
			if comparableErr.Equal(lm.list[k]) {
				return newIgnoreStrategy(err)
			}
		}
	}

	return nil
}

// ignoreStrategy holds data for an ignore strategy
type ignoreStrategy struct {
	previousErr error
}

// newIgnoreStrategy instantiates a new ignoreStrategy
func newIgnoreStrategy(previous error) ignoreStrategy {
	return ignoreStrategy{
		previousErr: previous,
	}
}

// previous holds the error that we wanted to ignore
func (as ignoreStrategy) previous() error {
	return as.previousErr
}

// last is defined to implement the interface
// returns nil since we are always mapping to nil for this strategy
func (as ignoreStrategy) last() error {
	return nil
}

// apply is defined to implement the interface
// returns nil since we are always mapping to nil for this strategy
func (as ignoreStrategy) apply() error {
	return nil
}
