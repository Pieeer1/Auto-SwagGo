package ext

func Distinct[T comparable](ts []T) []T {
	result := make([]T, 0)

	for _, tVal := range ts {
		if !Contains(result, tVal) {
			result = append(result, tVal)
		}
	}

	return result
}
