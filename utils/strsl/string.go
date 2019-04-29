package strsl

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// ContainsSub checks if baseSlice contains subSlice
func ContainsSub(baseSlice, subSlice []string) bool {
	if len(baseSlice) < len(subSlice) {
		return false
	}

	same := true

	for i := range subSlice {
		var found bool
		for j := range baseSlice {
			if subSlice[i] == baseSlice[j] {
				found = true
				break
			}
		}

		if !found {
			same = false
			break
		}
	}

	return same
}
