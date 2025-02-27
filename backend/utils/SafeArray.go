package utils

func GetSafeArray[T any](input []T) []T {
	if input == nil {
		// 返回一个空数组
		return []T{}
	}
	return input
}
