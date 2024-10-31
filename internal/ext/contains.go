package ext

func Contains[T comparable](ts []T, t T) bool {
	for _, tVal := range ts {
		if tVal == t {
			return true
		}
	}
	return false
}
