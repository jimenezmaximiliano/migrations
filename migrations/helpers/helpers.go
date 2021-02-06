package helpers

// AddTrailingSlashToPathIfNeeded takes a path and adds a trailing slash
// if there isn't any at the end of it.
func AddTrailingSlashToPathIfNeeded(path string) string {
	lastCharacterIndex := len(path) - 1
	if path[lastCharacterIndex:] != "/" {
		return path + "/"
	}

	return path
}
