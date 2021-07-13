package utils

func ContainsString(slice []string, needle string) bool {
	for _, n := range slice {
		if needle == n {
			return true
		}
	}
	return false
}

func PrependString(slice []string, e string) []string {
	slice = append(slice, "")
	copy(slice[1:], slice)
	slice[0] = e
	return slice
}
