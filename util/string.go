package util

func UniquifyStringSlice(slice []string) []string {
	if len(slice) < 2 {
		return slice
	}

	currItem := slice[0]
	for i := 1; i < len(slice); i++ {
		if slice[i] == currItem {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		} else {
			currItem = slice[i]
		}
	}
	return slice
}
