package utils

func DefaultValue[T any](old, new any) T {
	if new == "" || new == 0 {
		return old.(T)
	}
	return new.(T)
}
