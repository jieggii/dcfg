package util

func ItemIndex(array []string, item string) int {
	for i, x := range array {
		if x == item {
			return i
		}
	}
	return -1
}

func ItemIsInArray(array []string, item string) bool {
	if i := ItemIndex(array, item); i == -1 {
		return false
	}
	return true
}

func RemoveItem(array []string, index int) []string {
	return append(array[:index], array[index+1:]...)
}
