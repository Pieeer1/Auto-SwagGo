package ext

// SliceMap is a generic function that takes a slice of type T and a function that maps T to X and returns a slice of type X
// Usage:
// 	SliceMap([]int{1, 2, 3}, func(i int) string { return fmt.Sprintf("%d", i) }) // returns []string{"1", "2", "3"}
func SliceMap[T any, X any](ts []T, predicate func(T) X) []X {
	xs := make([]X, len(ts))
	for i, tVal := range ts {
		xs[i] = predicate(tVal)
	}
	return xs
}
