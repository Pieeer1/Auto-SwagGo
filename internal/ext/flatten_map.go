package ext

func FlattenMap[T any, X any](ts []T, predicate func(T) []X) []X {
	xs := make([]X, 0)
	for _, tVal := range ts {
		xs = append(xs, predicate(tVal)...)
	}
	return xs
}
