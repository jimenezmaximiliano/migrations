package helpers

func AddTrailingSlashToPathIfNeeded(path string) string {
	lastCharacterIndex := len(path) - 1
	if path[lastCharacterIndex:] != "/" {
		return path + "/"
	}

	return path
}
