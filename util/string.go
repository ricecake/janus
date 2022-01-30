package util

func UniquifyStringSlice(slice []string) []string {
	if len(slice) < 2 {
		return slice
	}

	currItem := slice[0]
	for i := 1; i < len(slice); i++ {
		// TODO: there's an optimization here for tracking the number of dupes,
		//       and then doing the append when we find a non-dupe
		//  if == { dupeCount++ } else { append(...); dupeCount=0; currItem = i }?  something like that.
		if slice[i] == currItem {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		} else {
			currItem = slice[i]
		}
	}
	return slice
}
