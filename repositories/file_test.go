package repositories_test

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jimenezmaximiliano/migrations/mocks"
	"github.com/jimenezmaximiliano/migrations/repositories"
)

func TestGettingMigrationFilePaths(test *testing.T) {
	test.Parallel()

	file := &mocks.File{}
	defer file.AssertExpectations(test)
	file.On("Name").Return("1_a.sql")
	file.On("IsDir").Return(false)
	files := []os.FileInfo{file}
	fileSystem := &mocks.FileSystem{}
	defer fileSystem.AssertExpectations(test)
	fileSystem.On("ReadDir", "/tmp/").Return(files, nil)
	repository := repositories.NewFileRepository(fileSystem)
	paths, err := repository.GetMigrationFilePaths("/tmp/")

	require.Nil(test, err)
	assert.Equal(test, "/tmp/1_a.sql", paths[0])
}

func TestGettingMigrationFilePathsFromADirectoryWithoutTrailingSlash(test *testing.T) {
	test.Parallel()

	file := &mocks.File{}
	defer file.AssertExpectations(test)
	file.On("Name").Return("1_a.sql")
	file.On("IsDir").Return(false)
	files := []os.FileInfo{file}
	fileSystem := &mocks.FileSystem{}
	defer fileSystem.AssertExpectations(test)
	fileSystem.On("ReadDir", "/tmp/").Return(files, nil)
	repository := repositories.NewFileRepository(fileSystem)
	paths, err := repository.GetMigrationFilePaths("/tmp")

	require.Nil(test, err)
	assert.Equal(test, "/tmp/1_a.sql", paths[0])
}

func TestGettingMigrationFilePathsOmitsNonSqlFilesAndDiretories(test *testing.T) {
	test.Parallel()

	file := &mocks.File{}
	defer file.AssertExpectations(test)
	file.On("Name").Return("1_a.sql")
	file.On("IsDir").Return(false)
	directory := &mocks.File{}
	defer directory.AssertExpectations(test)
	directory.On("Name").Return("2_b.sql")
	directory.On("IsDir").Return(true)
	txtFile := &mocks.File{}
	defer txtFile.AssertExpectations(test)
	txtFile.On("Name").Return("3.txt")
	txtFile.On("IsDir").Return(false)
	files := []os.FileInfo{directory, file, txtFile}
	fileSystem := &mocks.FileSystem{}
	defer fileSystem.AssertExpectations(test)
	fileSystem.On("ReadDir", "/tmp/").Return(files, nil)
	repository := repositories.NewFileRepository(fileSystem)
	paths, err := repository.GetMigrationFilePaths("/tmp")

	require.Nil(test, err)
	assert.Equal(test, "/tmp/1_a.sql", paths[0])
	assert.Len(test, paths, 1)
}

func TestGettingMigrationFilePathsFailsIfItIsNotPossibleToReadTheDirectory(test *testing.T) {
	test.Parallel()

	fileSystem := &mocks.FileSystem{}
	defer fileSystem.AssertExpectations(test)
	fileSystem.On("ReadDir", "/tmp/").Return(nil, fmt.Errorf("file system read error"))
	repository := repositories.NewFileRepository(fileSystem)
	paths, err := repository.GetMigrationFilePaths("/tmp/")

	assert.NotNil(test, err)
	assert.Nil(test, paths)
}

func TestGettingAQuery(test *testing.T) {
	test.Parallel()

	const query = "SELECT 1"
	fileSystem := &mocks.FileSystem{}
	defer fileSystem.AssertExpectations(test)
	fileSystem.On("ReadFile", "/tmp/1_a.sql").Return([]byte(query), nil)
	repository := repositories.NewFileRepository(fileSystem)
	readQuery, err := repository.GetMigrationQuery("/tmp/1_a.sql")

	assert.Equal(test, query, readQuery)
	assert.Nil(test, err)
}

func TestGettingAQueryFailsIfTheFileCannotBeRead(test *testing.T) {
	test.Parallel()

	fileSystem := &mocks.FileSystem{}
	defer fileSystem.AssertExpectations(test)
	fileSystem.On("ReadFile", "/tmp/1_a.sql").Return(nil, fmt.Errorf("file read error"))
	repository := repositories.NewFileRepository(fileSystem)
	readQuery, err := repository.GetMigrationQuery("/tmp/1_a.sql")

	assert.Equal(test, "", readQuery)
	assert.NotNil(test, err)
}

func TestCreatingAFile(test *testing.T) {
	test.Parallel()

	fileSystem := &mocks.FileSystem{}
	defer fileSystem.AssertExpectations(test)
	repository := repositories.NewFileRepository(fileSystem)

	fileSystem.
		On("WriteFile", "/tmp/1.sql", []byte("SELECT 1;"), fs.FileMode(0644)).
		Return(nil)

	err := repository.CreateMigration("/tmp/1.sql", "SELECT 1;")

	assert.Nil(test, err)
}
