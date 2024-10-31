package ext

func MapValues[T comparable, X any](ts map[T]X) []X {
	xs := make([]X, 0)
	for _, xVal := range ts {
		xs = append(xs, xVal)
	}
	return xs
}
