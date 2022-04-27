package utils

// Contains returns true if the target value is in the slice.
func Contains[T comparable](s []T, target T) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}
