package ext

import "slices"

// Push adds the element v to the beginning of the slice s.
func Push[T any](s []T, v T) []T {
	return slices.Insert(s, 0, v)
}
