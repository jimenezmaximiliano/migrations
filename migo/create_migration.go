package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/jimenezmaximiliano/migrations/helpers"
)

func createMigration(name, path string) (string, error) {
	pathWithTrailingSlash := helpers.AddTrailingSlashToPathIfNeeded(path)

	now := time.Now()
	fileName := fmt.Sprintf(
		"%d-%d-%d_%d-%d-%d_%s",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		name)

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
