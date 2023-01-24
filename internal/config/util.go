package config

func itemExists(array []string, item string) bool {
	for _, x := range array {
		if x == item {
			return true
		}
	}
	return false
}
