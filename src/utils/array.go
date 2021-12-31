package utils

func Contains(arr []string, e string) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}
