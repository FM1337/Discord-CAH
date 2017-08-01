package utils

// Fixes the index number for searching in lists
func IndexFixer(index int, listSize int) int {
	index = index - 1

	if index <= 0 {
		index = 0
	} else if index > listSize-1 {
		index = listSize - 1
	}

	return index
}
