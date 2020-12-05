package utils

// IsNameEqual is name equal
func IsNameEqual(r1 []string, r2 []string) bool {
	if len(r1) != len(r2) {
		return false
	}

	if len(r1) == 0 {
		return true
	}

	for idx, v := range r1 {
		if r2[idx] != v {
			return false
		}
	}

	return true
}
