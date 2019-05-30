package maperr

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

// Map an error to another error
func (hm HashableMapper) Map(err error) MapResult {
	var toMap = err
	previous := LastAppended(err)
	if previous != nil {
		toMap = previous
	}
	key := hm.tryMakeHashable(toMap)
	mapped, ok := hm[key]
	if !ok {
		return nil
	}
	return NewAppendStrategy(err, mapped)
}

func (hm HashableMapper) tryMakeHashable(err error) error {
	key := err
	if casted, ok := err.(formattedError); ok {
		key = casted.Hashable()
	}
	return key
}
