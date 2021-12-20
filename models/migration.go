package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	// StatusUnknown represents a MigrationContainer without a set status (default value).
	StatusUnknown int8 = 0
	// StatusNotRun represents a MigrationContainer that hasn't been run yet.
	StatusNotRun int8 = 1
	// StatusSuccessful represents a MigrationContainer that has been run and it was successful.
	StatusSuccessful int8 = 2
	// StatusFailed represents a MigrationContainer that has been run and it failed.
	StatusFailed int8 = -1
)

// Migration represents a database MigrationContainer and its state (immutable).
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

type MigrationContainer struct {
	absolutePath string
	name         string
	status       int8
	query        string
	err          error
	order        uint64
}

// Ensure MigrationContainer implements Migration
var _ Migration = MigrationContainer{}

// GetAbsolutePath returns the absolute path of the MigrationContainer file.
func (thisMigration MigrationContainer) GetAbsolutePath() string {
	return thisMigration.absolutePath
}

// GetName returns the file name of the MigrationContainer file.
func (thisMigration MigrationContainer) GetName() string {
	return thisMigration.name
}

// GetStatus returns the current status of the MigrationContainer using the constants on this package.
func (thisMigration MigrationContainer) GetStatus() int8 {
	return thisMigration.status
}

// GetOrder returns the order on which the MigrationContainer should be run.
func (thisMigration MigrationContainer) GetOrder() uint64 {
	return thisMigration.order
}

// ShouldBeRun returns true if the MigrationContainer has not been run yet.
func (thisMigration MigrationContainer) ShouldBeRun() bool {
	return thisMigration.status == StatusNotRun
}

// WasSuccessful returns true if the current status is StatusSuccessful.
func (thisMigration MigrationContainer) WasSuccessful() bool {
	return thisMigration.status == StatusSuccessful
}

// HasFailed returns true if the current status is StatusFailed.
func (thisMigration MigrationContainer) HasFailed() bool {
	return thisMigration.status == StatusFailed
}

// GetQuery returns the sql query of the MigrationContainer.
func (thisMigration MigrationContainer) GetQuery() string {
	return thisMigration.query
}

// NewAsFailed returns a copy of the MigrationContainer but with a StatusFailed status.
func (thisMigration MigrationContainer) NewAsFailed(err error) Migration {
	return MigrationContainer{
		absolutePath: thisMigration.absolutePath,
		name:         thisMigration.name,
		status:       StatusFailed,
		query:        thisMigration.query,
		err:          err,
		order:        thisMigration.order,
	}
}

// NewAsNotRun returns a copy of the MigrationContainer but with a StatusNotRun status.
func (thisMigration MigrationContainer) NewAsNotRun() Migration {
	newMigration, _ := NewMigration(thisMigration.GetAbsolutePath(), thisMigration.GetQuery(), StatusNotRun)

	return newMigration
}

// NewAsSuccessful returns a copy of the MigrationContainer but with a StatusSuccessful status.
func (thisMigration MigrationContainer) NewAsSuccessful() Migration {
	newMigration, _ := NewMigration(thisMigration.GetAbsolutePath(), thisMigration.GetQuery(), StatusSuccessful)

	return newMigration
}

// ShouldBeRunFirst returns true if this MigrationContainer needs to be run before the given MigrationContainer
// (used for sorting migrations).
func (thisMigration MigrationContainer) ShouldBeRunFirst(anotherMigration Migration) bool {
	return thisMigration.GetOrder() < anotherMigration.GetOrder()
}

// NewMigration is a constructor for a Migration implementation.
func NewMigration(absolutePath string, query string, status int8) (Migration, error) {
	if status < -1 || status > 2 {
		return MigrationContainer{}, fmt.Errorf("MigrationContainer invalid status (status: %d)", status)
	}

	fileName := extractFileName(absolutePath)
	order, err := getOrderFromFileName(fileName)
	if err != nil {
		return MigrationContainer{}, err
	}

	return MigrationContainer{
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
		return 0, errors.Wrapf(err, "invalid MigrationContainer file name [%s]", orderAsString)
	}

	return order, nil
}

// GetError returns the error that caused the MigrationContainer to fail.
func (thisMigration MigrationContainer) GetError() error {
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
