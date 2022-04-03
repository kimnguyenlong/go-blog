package utils

func ArrayContains[E comparable](arr []E, e E) bool {
	for _, v := range arr {
		if v == e {
			return true
		}
	}
	return false
}
