package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/jimenezmaximiliano/migrations/helpers"
)

// Deprecated: migo is deprecated.
func createMigration(name, path string) (string, error) {
	if name == "" {
		return "", errors.New("a migration's name cannot be empty")
	}

	pathWithTrailingSlash := helpers.AddTrailingSlashToPathIfNeeded(path)

	now := time.Now()
	fileName := fmt.Sprintf("%d_%s", now.UnixNano(), name)

	filePath := pathWithTrailingSlash + fileName
	extension := filePath[len(filePath)-4:]
	if extension != ".sql" {
		filePath += ".sql"
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create migration file on %s", filePath)
	}
	defer file.Close()

	_, err = file.WriteString("--\nSELECT 1;")
	if err != nil {
		return "", errors.Wrapf(err, "failed to write content to migration file on %s", filePath)
	}

	return filePath, nil
}
