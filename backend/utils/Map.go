package utils

func Map[T, U any](s []T, f func(T) U) []U {
	result := make([]U, len(s))
	for i := range s {
		result[i] = f(s[i])
	}
	return result
}
