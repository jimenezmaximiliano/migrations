package repositories

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/jimenezmaximiliano/migrations/mocks"
)

func TestCreatingTheMigrationsTable(test *testing.T) {
	db := &mocks.DB{}
	db.On("Exec", mock.AnythingOfType("string")).Return(nil, nil)
	repository := NewDBRepository(db)
	err := repository.CreateMigrationsTableIfNeeded()

	if err != nil {
		test.Error(err)
	}
}

func TestCreatingTheMigrationsTableFailsIfThereWasAnError(test *testing.T) {
	db := &mocks.DB{}
	db.On("Exec", mock.AnythingOfType("string")).Return(nil, fmt.Errorf("db exec error"))
	repository := NewDBRepository(db)
	err := repository.CreateMigrationsTableIfNeeded()

	if err == nil {
		test.Fail()
	}
}

func TestPingingAnOkConnection(test *testing.T) {
	db := &mocks.DB{}
	db.On("Ping").Return(nil)
	repository := NewDBRepository(db)
	err := repository.Ping()

	if err != nil {
		test.Fail()
	}
}

func TestPingingAKOConnection(test *testing.T) {
	db := &mocks.DB{}
	db.On("Ping").Return(fmt.Errorf("db ping error"))
	repository := NewDBRepository(db)
	err := repository.Ping()

	if err == nil {
		test.Fail()
	}
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

	if err != nil {
		test.Error(err)
	}

	if filePaths[0] != "/tmp/migrationAlreadyRun.sql" {
		test.Errorf("Expected migration path /tmp/migrationAlreadyRun.sql but got %s", filePaths[0])
	}
}

func TestGettingAlreadyRunMigrationFilePathsFailsIfTheQueryFails(test *testing.T) {
	db := &mocks.DB{}
	db.On("Query", mock.AnythingOfType("string")).Return(nil, fmt.Errorf("db query error"))
	repository := NewDBRepository(db)
	filePaths, err := repository.GetAlreadyRunMigrationFilePaths("/tmp/")

	if err == nil || filePaths != nil {
		test.Fail()
	}
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

	if err == nil || filePaths != nil {
		test.Fail()
	}
}

func TestRunningASuccessfulMigrationQuery(test *testing.T) {
	const query = "SELECT 1"
	db := &mocks.DB{}
	db.On("Exec", query).Return(nil, nil)
	repository := NewDBRepository(db)
	err := repository.RunMigrationQuery(query)

	if err != nil {
		test.Fail()
	}
}

func TestRunningABrokenMigrationQueryFails(test *testing.T) {
	const query = "SELECT * FROM"
	db := &mocks.DB{}
	db.On("Exec", query).Return(nil, fmt.Errorf("db query error"))
	repository := NewDBRepository(db)
	err := repository.RunMigrationQuery(query)

	if err == nil {
		test.Fail()
	}
}

func TestRegisteringARunMigration(test *testing.T) {
	const migrationName = ""
	db := &mocks.DB{}
	db.On("Exec", mock.AnythingOfType("string"), migrationName).Return(nil, nil)
	repository := NewDBRepository(db)
	err := repository.RegisterRunMigration(migrationName)

	if err != nil {
		test.Fail()
	}
}

func TestRegisteringARunMigrationFailsIfTheInsertFails(test *testing.T) {
	const migrationName = ""
	db := &mocks.DB{}
	db.On("Exec", mock.AnythingOfType("string"), migrationName).Return(nil, fmt.Errorf("db query error"))
	repository := NewDBRepository(db)
	err := repository.RegisterRunMigration(migrationName)

	if err == nil {
		test.Fail()
	}
}