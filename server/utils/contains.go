package utils

func Contains(slice []string, item string) (int, bool) {
	for i, v := range slice {
		if v == item {
			return i, true
		}
	}
	return -1, false
}
