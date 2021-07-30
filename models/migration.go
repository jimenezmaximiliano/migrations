package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	// StatusUnknown represents a migration without a set status (default value).
	StatusUnknown int8 = 0
	// StatusNotRun represents a migration that hasn't been run yet.
	StatusNotRun int8 = 1
	// StatusSuccessful represents a migration that has been run and it was successful.
	StatusSuccessful int8 = 2
	// StatusFailed represents a migration that has been run and it failed.
	StatusFailed int8 = -1
)

// Migration represents a database migration and its state (immutable).
type Migration interface {
	GetAbsolutePath() string
	GetName() string
	GetStatus() int8
	GetOrder() uint64
	ShouldBeRun() bool
	GetQuery() string
	NewAsFailed(err error) Migration
	NewAsSuccessful() Migration
	NewAsNotRun() Migration
	WasSuccessful() bool
	HasFailed() bool
	ShouldBeRunFirst(anotherMigration Migration) bool
	GetError() error
}

type migration struct {
	absolutePath string
	name         string
	status       int8
	query        string
	err			 error
	order        uint64
}

// Ensure migration implements Migration
var _ Migration = migration{}

// GetAbsolutePath returns the absolute path of the migration file.
func (thisMigration migration) GetAbsolutePath() string {
	return thisMigration.absolutePath
}

// GetName returns the file name of the migration file.
func (thisMigration migration) GetName() string {
	return thisMigration.name
}

// GetStatus returns the current status of the migration using the constants on this package.
func (thisMigration migration) GetStatus() int8 {
	return thisMigration.status
}

// GetOrder returns the order on which the migration should be run.
func (thisMigration migration) GetOrder() uint64 {
	return thisMigration.order
}

// ShouldBeRun returns true if the migration has not been run yet.
func (thisMigration migration) ShouldBeRun() bool {
	return thisMigration.status == StatusNotRun
}

// WasSuccessful returns true if the current status is StatusSuccessful.
func (thisMigration migration) WasSuccessful() bool {
	return thisMigration.status == StatusSuccessful
}

// HasFailed returns true if the current status is StatusFailed.
func (thisMigration migration) HasFailed() bool {
	return thisMigration.status == StatusFailed
}

// GetQuery returns the sql query of the migration.
func (thisMigration migration) GetQuery() string {
	return thisMigration.query
}

// NewAsFailed returns a copy of the migration but with a StatusFailed status.
func (thisMigration migration) NewAsFailed(err error) Migration {
	return migration{
		absolutePath: thisMigration.absolutePath,
		name:         thisMigration.name,
		status:       StatusFailed,
		query:        thisMigration.query,
		err:          err,
		order:        thisMigration.order,
	}
}

// NewAsNotRun returns a copy of the migration but with a StatusNotRun status.
func (thisMigration migration) NewAsNotRun() Migration {
	newMigration, _ := NewMigration(thisMigration.GetAbsolutePath(), thisMigration.GetQuery(), StatusNotRun)

	return newMigration
}

// NewAsSuccessful returns a copy of the migration but with a StatusSuccessful status.
func (thisMigration migration) NewAsSuccessful() Migration {
	newMigration, _ := NewMigration(thisMigration.GetAbsolutePath(), thisMigration.GetQuery(), StatusSuccessful)

	return newMigration
}

// ShouldBeRunFirst returns true if this migration needs to be run before the given migration
// (used for sorting migrations).
func (thisMigration migration) ShouldBeRunFirst(anotherMigration Migration) bool {
	return thisMigration.GetOrder() < anotherMigration.GetOrder()
}

// NewMigration is a constructor for a Migration implementation.
func NewMigration(absolutePath string, query string, status int8) (Migration, error) {
	if status < -1 || status > 2 {
		return migration{}, fmt.Errorf("migration invalid status (status: %d)", status)
	}

	fileName := extractFileName(absolutePath)
	order, err := getOrderFromFileName(fileName)
	if err != nil {
		return migration{}, err
	}

	return migration{
		absolutePath: absolutePath,
		name:         fileName,
		status:       status,
		query:        query,
		order:        order,
	}, nil
}

func getOrderFromFileName(fileName string) (uint64, error) {
	result := strings.Split(fileName, "_")
	orderAsString := result[0]
	order, err := strconv.ParseUint(orderAsString, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid migration file name [%s]", orderAsString)
	}

	return order, nil
}

// GetError returns the error that caused the migration to fail.
func (thisMigration migration) GetError() error {
	return thisMigration.err
}

func extractFileName(absolutePath string) string {
	absolutePathParts := strings.Split(absolutePath, "/")

	return getSliceLastElement(absolutePathParts)
}

func getSliceLastElement(theSlice []string) string {
	lastIndex := len(theSlice) - 1

	return theSlice[lastIndex]
}
