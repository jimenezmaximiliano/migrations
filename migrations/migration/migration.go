package migration

import (
	"fmt"
	"strings"
)

const (
	StatusUnknown = 0
	StatusNotRun = 1
	StatusSuccessful = 2
	StatusFailed = -1
)

type Migration interface {
	getAbsolutePath() string
	getName() string
	getStatus() int8
}

type migration struct {
	absolutePath string
	name string
	status int8
}

func (migration migration) getAbsolutePath() string {
	return migration.absolutePath
}

func (migration migration) getName() string {
	return migration.name
}

func (migration migration) getStatus() int8 {
	return migration.status
}

func New(absolutePath string, status int8) (Migration, error) {

	if status < -1 || status > 2 {
		return migration{}, fmt.Errorf("verySimpleMigrations.migration.New.invalidStatus (status: %d)", status)
	}

	return migration{
		absolutePath: absolutePath,
		name: extractFileName(absolutePath),
		status: status,
	}, nil
}

func extractFileName(absolutePath string) string {
	parts := strings.Split(absolutePath, "/")

	return parts[len(parts)-1]
}