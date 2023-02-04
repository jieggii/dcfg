package config

func itemIndex(array []string, item string) int {
	for i, x := range array {
		if x == item {
			return i
		}
	}
	return -1
}

func itemIsInArray(array []string, item string) bool {
	if i := itemIndex(array, item); i == -1 {
		return false
	}
	return true
}

func removeItem(array []string, index int) []string {
	return append(array[:index], array[index+1:]...)
}
