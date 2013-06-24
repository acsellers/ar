package ar

func stringMatch(items []string, wanted string) bool {
	for _, item := range items {
		if item == wanted {
			return true
		}
	}

	return false
}
