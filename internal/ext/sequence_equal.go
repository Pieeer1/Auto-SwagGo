package ext

// SequenceEqual compares two slices of comparable types and returns true if they are equal.
func SequenceEqual[T comparable](ts1 []T, ts2 []T) bool {
	if len(ts1) != len(ts2) {
		return false
	}
	for i, t1 := range ts1 {
		if t1 != ts2[i] {
			return false
		}
	}
	return true
}
