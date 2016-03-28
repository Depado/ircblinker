package utils

// StringInSlice checks if a string is present in an array of strings.
func StringInSlice(a string, l []string) bool {
	for _, b := range l {
		if b == a {
			return true
		}
	}
	return false
}

// Any checks whether any string of first argument is in the second one.
func Any(a, b []string) bool {
	for _, t := range a {
		if StringInSlice(t, b) {
			return true
		}
	}
	return false
}
