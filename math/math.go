package math

import "golang.org/x/exp/constraints"

// Max returns max number.
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}

	return b
}
