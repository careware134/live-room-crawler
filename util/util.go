package util

func isInStringList(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

func IsInList(num int, list []int) bool {
	for _, s := range list {
		if s == num {
			return true
		}
	}
	return false
}
