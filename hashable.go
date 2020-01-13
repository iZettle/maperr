package maperr

import (
	"errors"

	"go.uber.org/multierr"
)

// HashableMapper simple implementation of Mapper which only works
// on hashable error keys
type HashableMapper map[error]error

// NewHashableMapper make a new instance of HashableMapper
func NewHashableMapper() HashableMapper {
	return HashableMapper{}
}

// Append append an error to error association
func (hm HashableMapper) Append(err, match error) HashableMapper {
	key := hm.tryMakeHashable(err)
	val := hm.tryMakeHashable(match)
	hm[key] = val
	return hm
}

// mapErr an error to another error
func (hm HashableMapper) mapErr(err error) mapResult {
	errorsToMap := []error{
		err,
	}
	if errList := multierr.Errors(err); len(errList) > 0 {
		errorsToMap = errList
	}

	for i := len(errorsToMap) - 1; i >= 0; i-- {
		key := hm.tryMakeHashable(errorsToMap[i])
		mapped, ok := hm[key]
		if ok {
			return newAppendStrategy(err, mapped)
		}
	}
	return nil
}

func (hm HashableMapper) tryMakeHashable(err error) error {
	key := err

	var ferr formattedError
	if errors.As(err, &ferr) {
		key = ferr.Hashable()
	}

	return key
}
