package repositories

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jimenezmaximiliano/migrations/mocks"
)

func TestCreatingTheMigrationsTable(test *testing.T) {
	db := &mocks.DB{}
	db.On("Exec", mock.AnythingOfType("string")).Return(nil, nil)
	repository := NewDBRepository(db)
	err := repository.CreateMigrationsTableIfNeeded()

	assert.Nil(test, err)
}

func TestCreatingTheMigrationsTableFailsIfThereWasAnError(test *testing.T) {
	db := &mocks.DB{}
	db.On("Exec", mock.AnythingOfType("string")).Return(nil, fmt.Errorf("db exec error"))
	repository := NewDBRepository(db)
	err := repository.CreateMigrationsTableIfNeeded()

	assert.NotNil(test, err)
}

func TestPingingAnOkConnection(test *testing.T) {
	db := &mocks.DB{}
	db.On("Ping").Return(nil)
	repository := NewDBRepository(db)
	err := repository.Ping()

	assert.Nil(test, err)
}

func TestPingingAKOConnection(test *testing.T) {
	db := &mocks.DB{}
	db.On("Ping").Return(fmt.Errorf("db ping error"))
	repository := NewDBRepository(db)
	err := repository.Ping()

	assert.NotNil(test, err)
}

func TestGettingAlreadyRunMigrationFilePaths(test *testing.T) {
	rows := &mocks.DBRows{}
	rows.On("Close").Return(nil).Once()
	rows.On("Next").Return(true).Once()
	rows.On("Next").Return(false).Once()
	rows.On("Scan", mock.AnythingOfType("*string")).Return(nil).Run(func(args mock.Arguments) {
		var thePath *string = args[0].(*string)
		*thePath = "migrationAlreadyRun.sql"
	})
	db := &mocks.DB{}
	db.On("Query", mock.AnythingOfType("string")).Return(rows, nil)
	repository := NewDBRepository(db)
	filePaths, err := repository.GetAlreadyRunMigrationFilePaths("/tmp/")

	require.Nil(test, err)
	assert.Equal(test, "/tmp/migrationAlreadyRun.sql", filePaths[0])
}

func TestGettingAlreadyRunMigrationFilePathsFailsIfTheQueryFails(test *testing.T) {
	db := &mocks.DB{}
	db.On("Query", mock.AnythingOfType("string")).Return(nil, fmt.Errorf("db query error"))
	repository := NewDBRepository(db)
	filePaths, err := repository.GetAlreadyRunMigrationFilePaths("/tmp/")

	assert.NotNil(test, err)
	assert.Nil(test, filePaths)
}

func TestGettingAlreadyRunMigrationFilePathsFailsIfRowsCannotBeScanned(test *testing.T) {
	rows := &mocks.DBRows{}
	rows.On("Close").Return(nil).Once()
	rows.On("Next").Return(true).Once()
	rows.On("Next").Return(false).Once()
	rows.On("Scan", mock.AnythingOfType("*string")).Return(fmt.Errorf("rows scan error"))
	db := &mocks.DB{}
	db.On("Query", mock.AnythingOfType("string")).Return(rows, nil)
	repository := NewDBRepository(db)
	filePaths, err := repository.GetAlreadyRunMigrationFilePaths("/tmp/")

	assert.NotNil(test, err)
	assert.Nil(test, filePaths)
}

func TestRunningASuccessfulMigrationQuery(test *testing.T) {
	const query = "SELECT 1"
	db := &mocks.DB{}
	db.On("Exec", query).Return(nil, nil)
	repository := NewDBRepository(db)
	err := repository.RunMigrationQuery(query)

	assert.Nil(test, err)
}

func TestRunningABrokenMigrationQueryFails(test *testing.T) {
	const query = "SELECT * FROM"
	db := &mocks.DB{}
	db.On("Exec", query).Return(nil, fmt.Errorf("db query error"))
	repository := NewDBRepository(db)
	err := repository.RunMigrationQuery(query)

	assert.NotNil(test, err)
}

func TestRegisteringARunMigration(test *testing.T) {
	const migrationName = ""
	db := &mocks.DB{}
	db.On("Exec", mock.AnythingOfType("string"), migrationName).Return(nil, nil)
	repository := NewDBRepository(db)
	err := repository.RegisterRunMigration(migrationName)

	assert.Nil(test, err)
}

func TestRegisteringARunMigrationFailsIfTheInsertFails(test *testing.T) {
	const migrationName = ""
	db := &mocks.DB{}
	db.On("Exec", mock.AnythingOfType("string"), migrationName).Return(nil, fmt.Errorf("db query error"))
	repository := NewDBRepository(db)
	err := repository.RegisterRunMigration(migrationName)

	assert.NotNil(test, err)
}