package maperr

// IgnoreListMapper is a mapper that allow to specify a list of error that we
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
func (lm IgnoreListMapper) Append(err Error) IgnoreListMapper {
	lm.list = append(lm.list, err)
	return lm
}

// Map an error to an ignore strategy
func (lm IgnoreListMapper) Map(err error) MapResult {
	var toMap = err
	previous := LastAppended(err)
	if previous != nil {
		toMap = previous
	}

	comparableErr, ok := toMap.(Error)
	if !ok {
		comparableErr = NewError(toMap.Error())
	}

	for k := range lm.list {
		if !comparableErr.Equal(lm.list[k]) {
			continue
		}
		return NewIgnoreStrategy(err)
	}
	return nil
}

// IgnoreStrategy holds data for an ignore strategy
type IgnoreStrategy struct {
	previous error
}

// NewIgnoreStrategy instantiates a new IgnoreStrategy
func NewIgnoreStrategy(previous error) IgnoreStrategy {
	return IgnoreStrategy{
		previous: previous,
	}
}

// Previous holds the error that we wanted to ignore
func (as IgnoreStrategy) Previous() error {
	return as.previous
}

// Last is defined to implement the interface
// returns nil since we are always mapping to nil for this strategy
func (as IgnoreStrategy) Last() error {
	return nil
}

// Apply is defined to implement the interface
// returns nil since we are always mapping to nil for this strategy
func (as IgnoreStrategy) Apply() error {
	return nil
}
