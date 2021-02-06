package migration

import (
	"fmt"
	"sort"
	"strings"
)

const (
	// StatusUnknown represents a migration without a set status (default value).
	StatusUnknown = 0
	// StatusNotRun represents a migration that hasn't been run yet.
	StatusNotRun = 1
	// StatusSuccessful represents a migration that has been run and it was successful.
	StatusSuccessful = 2
	// StatusFailed represents a migration that has been run and it failed.
	StatusFailed = -1
)

// Migration represents a database migration and its state (immutable).
type Migration interface {
	GetAbsolutePath() string
	GetName() string
	GetStatus() int8
	ShouldBeRun() bool
	GetQuery() string
	NewAsFailed() Migration
	NewAsSuccessful() Migration
	WasSuccessful() bool
	HasFailed() bool
	ShouldBeRunFirst(anotherMigration Migration) bool
}

type migration struct {
	absolutePath string
	name         string
	status       int8
	query        string
}

// GetAbsolutePath returns the absolute path of the migration file.
func (migration migration) GetAbsolutePath() string {
	return migration.absolutePath
}

// GeName returns the file name of the migration file.
func (migration migration) GetName() string {
	return migration.name
}

// GetStatus returns the current status of the migration using the consts on this package.
func (migration migration) GetStatus() int8 {
	return migration.status
}

// ShouldBeRun returns true if the migration has not been run yet.
func (migration migration) ShouldBeRun() bool {
	return migration.status != StatusSuccessful
}

// WasSuccessful returs true if the current status is StatusSuccessful.
func (migration migration) WasSuccessful() bool {
	return migration.status == StatusSuccessful
}

// HasFailed returns true if the current status is StatusFailed.
func (migration migration) HasFailed() bool {
	return migration.status == StatusFailed
}

// GetQuery returns the sql query of the migration.
func (migration migration) GetQuery() string {
	return migration.query
}

// NewAsFailed returns a copy of the migration but with a StatusFailed status.
func (migration migration) NewAsFailed() Migration {
	newMigration, _ := NewMigration(migration.GetAbsolutePath(), migration.GetQuery(), StatusFailed)

	return newMigration
}

// NewAsFailed returns a copy of the migration but with a StatusSuccessful status.
func (migration migration) NewAsSuccessful() Migration {
	newMigration, _ := NewMigration(migration.GetAbsolutePath(), migration.GetQuery(), StatusSuccessful)

	return newMigration
}

// ShouldBeRunFirst returns true if this migration needs to be run before the given migration
// (used for sorting migrations).
func (migration migration) ShouldBeRunFirst(anotherMigration Migration) bool {
	names := []string{
		migration.name,
		anotherMigration.GetName(),
	}
	sort.Strings(names)

	return names[0] == migration.name
}

// NewMigration is a constructor for a Migration implementation.
func NewMigration(absolutePath string, query string, status int8) (Migration, error) {
	if status < -1 || status > 2 {
		return migration{}, fmt.Errorf("migration invalid status (status: %d)", status)
	}

	return migration{
		absolutePath: absolutePath,
		name:         extractFileName(absolutePath),
		status:       status,
		query:        query,
	}, nil
}

func extractFileName(absolutePath string) string {
	absolutePathParts := strings.Split(absolutePath, "/")

	return getSliceLastElement(absolutePathParts)
}

func getSliceLastElement(theSlice []string) string {
	lastIndex := len(theSlice) - 1

	return theSlice[lastIndex]
}
