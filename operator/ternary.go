package operator

// Ternary returns a if condition is true, otherwise it returns b.
func Ternary[T comparable](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}
