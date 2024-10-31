package ext

func Where[T any](ts []T, predicate func(t T) bool) []T {
	result := make([]T, 0)

	for _, tVal := range ts {
		if predicate(tVal) {
			result = append(result, tVal)
		}
	}

	return result
}
