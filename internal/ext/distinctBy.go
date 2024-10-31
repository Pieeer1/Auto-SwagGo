package ext

func DistinctBy[T any, X comparable](ts []T, predicate func(T) X) []T {
	visited := make([]X, 0)
	result := make([]T, 0)

	for _, tVal := range ts {
		x := predicate(tVal)
		if !Contains(visited, x) {
			visited = append(visited, x)
			result = append(result, tVal)
		}
	}

	return result
}
