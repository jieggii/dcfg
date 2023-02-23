package util

// ItemIndex returns index of item in string array if item is present.
// Otherwise it returns -1.
func ItemIndex(array []string, item string) int {
	for i, x := range array {
		if x == item {
			return i
		}
	}
	return -1
}

// ItemIsInArray returns true if item is in a string array.
func ItemIsInArray(array []string, item string) bool {
	if i := ItemIndex(array, item); i == -1 {
		return false
	}
	return true
}

// RemoveItem removes item from string array.
func RemoveItem(array []string, index int) []string {
	return append(array[:index], array[index+1:]...)
}
