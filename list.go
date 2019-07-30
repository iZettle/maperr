package maperr

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
func (lm ListMapper) Appendf(format string, match Error) ListMapper {
	return lm.Append(Errorf(format), match)
}

// Append append an error to error association
func (lm ListMapper) Append(err, match Error) ListMapper {
	lm.errorPairs = append(lm.errorPairs,
		PairErrors{
			err:   err,
			match: match,
		})
	return lm
}

// Map a formatted error to an error
func (lm ListMapper) Map(err error) MapResult {
	var toMap = err
	previous := LastAppended(err)
	if previous != nil {
		toMap = previous
	}

	comparableErr := CastError(toMap)
	for k := range lm.errorPairs {
		if !comparableErr.Equal(lm.errorPairs[k].err) {
			continue
		}
		return NewAppendStrategy(err, lm.errorPairs[k].match)
	}
	return nil
}

// AppendStrategy holds data for an ignore strategy
type AppendStrategy struct {
	previous error
	last     error
}

// NewAppendStrategy instantiates a new AppendStrategy
func NewAppendStrategy(previous, last error) AppendStrategy {
	return AppendStrategy{previous: previous, last: last}
}

// Previous returns the error that we want to append to
func (as AppendStrategy) Previous() error {
	return as.previous
}

// Last returns the error that we are appending
func (as AppendStrategy) Last() error {
	return as.last
}

// Apply the append strategy by appending previous to last
func (as AppendStrategy) Apply() error {
	if as.last == nil {
		return nil
	}
	return Append(as.previous, as.last)
}
