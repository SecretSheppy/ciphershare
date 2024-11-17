package utils

// StringInList returns true if the string is present in the list.
func StringInList(a string, list []string) bool {
	m := make(map[string]bool)
	for _, item := range list {
		m[item] = true
	}

	return m[a]
}
