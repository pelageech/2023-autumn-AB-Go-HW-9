package iterator

// Interface defines a default iterator. It is used for
// Iterate function. You should use this function for iterating.
type Interface[T any] interface {
	// HasNext defines if the next iteration exist.
	HasNext() bool

	// Next returns the next thing in iteration or error
	// in some unexpected situations.
	Next() (T, error)
}

// Iterate iterates over Interface implementation and calls f
// on each iteration. If i.Next() of f returns error, the iteration
// stops and the error that is returned by functions above is returned.
func Iterate[T any](i Interface[T], f func(t T) error) error {
	if !i.HasNext() {
		return nil
	}

	for t, err := i.Next(); i.HasNext(); t, err = i.Next() {
		if err != nil {
			return err
		}

		if err := f(t); err != nil {
			return err
		}
	}

	return nil
}
