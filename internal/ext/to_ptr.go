package ext

func ToPtr[T any](v T) *T {
	return &v
}
