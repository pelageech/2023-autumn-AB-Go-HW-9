package iterator

import "fmt"

// Simple defines a default iterator. It is used for
// Iterate function. You should use this function for iterating.
type Simple[T any] interface {
	// HasNext defines if the next iteration exists.
	HasNext() bool

	// Next returns the next thing in iteration or error
	// in some unexpected situations.
	Next() (T, error)
}

// Iterate iterates over Simple implementation and calls f
// on each iteration. If i.Next() of f returns error, the iteration
// stops and the error returned by the functions above is returned.
func Iterate[T any](i Simple[T], f func(t T) error) error {
	if !i.HasNext() {
		return nil
	}

	for t, err := i.Next(); i.HasNext(); t, err = i.Next() {
		if err != nil {
			return fmt.Errorf("taking next value error: %w", err)
		}

		if err := f(t); err != nil {
			return fmt.Errorf("parameterized function error: %w", err)
		}
	}

	return nil
}
