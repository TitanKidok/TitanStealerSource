package endpoint

func SliceCheck(a string, list []string) (bool, string) {
	for _, b := range list {
		if b == a {
			return true, b
		}
	}
	return false, ""
}
