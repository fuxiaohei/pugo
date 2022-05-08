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

// UniqueStringsSlice returns a unique slice of strings.
func UniqueStringsSlice(strSlice []string) []string {
	keys := make(map[string]struct{})
	list := []string{}

	for _, entry := range strSlice {
		if _, ok := keys[entry]; !ok {
			keys[entry] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}
