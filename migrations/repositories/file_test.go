package repositories

import (
	"fmt"
	"os"
	"testing"

	"github.com/jimenezmaximiliano/migrations/mocks"
)

func TestGettingMigrationFilePaths(test *testing.T) {
	file := &mocks.File{}
	file.On("Name").Return("1.sql")
	file.On("IsDir").Return(false)
	files := []os.FileInfo{file}
	fileSystem := &mocks.FileSystem{}
	fileSystem.On("ReadDir", "/tmp/").Return(files, nil)
	repository := NewFileRepository(fileSystem)
	paths, err := repository.GetMigrationFilePaths("/tmp/")

	if err != nil {
		test.Error(err)
	}

	if paths[0] != "/tmp/1.sql" {
		test.Errorf("Expected /tmp/1.sql but got %s", paths[0])
	}
}

func TestGettingMigrationFilePathsFromADirectoryWithoutTrailingSlash(test *testing.T) {
	file := &mocks.File{}
	file.On("Name").Return("1.sql")
	file.On("IsDir").Return(false)
	files := []os.FileInfo{file}
	fileSystem := &mocks.FileSystem{}
	fileSystem.On("ReadDir", "/tmp/").Return(files, nil)
	repository := NewFileRepository(fileSystem)
	paths, err := repository.GetMigrationFilePaths("/tmp")

	if err != nil {
		test.Error(err)
	}

	if paths[0] != "/tmp/1.sql" {
		test.Errorf("Expected /tmp/1.sql but got %s", paths[0])
	}
}

func TestGettingMigrationFilePathsOmitsNonSqlFilesAndDiretories(test *testing.T) {
	file := &mocks.File{}
	file.On("Name").Return("1.sql")
	file.On("IsDir").Return(false)
	directory := &mocks.File{}
	directory.On("Name").Return("2.sql")
	directory.On("IsDir").Return(true)
	txtFile := &mocks.File{}
	txtFile.On("Name").Return("3.txt")
	txtFile.On("IsDir").Return(false)
	files := []os.FileInfo{directory, file, txtFile}
	fileSystem := &mocks.FileSystem{}
	fileSystem.On("ReadDir", "/tmp/").Return(files, nil)
	repository := NewFileRepository(fileSystem)
	paths, err := repository.GetMigrationFilePaths("/tmp")

	if err != nil {
		test.Error(err)
	}

	if paths[0] != "/tmp/1.sql" {
		test.Errorf("Expected /tmp/1.sql but got %s", paths[0])
	}

	if len(paths) != 1 {
		test.Fail()
	}
}

func TestGettingMigrationFilePathsFailsIfItIsNotPossibleToReadTheDirectory(test *testing.T) {
	fileSystem := &mocks.FileSystem{}
	fileSystem.On("ReadDir", "/tmp/").Return(nil, fmt.Errorf("file system read error"))
	repository := NewFileRepository(fileSystem)
	paths, err := repository.GetMigrationFilePaths("/tmp/")

	if err == nil || paths != nil {
		test.Error(err)
	}
}

func TestGettingAQuery(test *testing.T) {
	const query = "SELECT 1"
	fileSystem := &mocks.FileSystem{}
	fileSystem.On("ReadFile", "/tmp/1.sql").Return([]byte(query), nil)
	repository := NewFileRepository(fileSystem)
	readQuery, err := repository.GetMigrationQuery("/tmp/1.sql")

	if readQuery != query || err != nil {
		test.Fail()
	}
}

func TestGettingAQueryFailsIfTheFileCannotBeRead(test *testing.T) {
	const query = "SELECT 1"
	fileSystem := &mocks.FileSystem{}
	fileSystem.On("ReadFile", "/tmp/1.sql").Return(nil, fmt.Errorf("file read error"))
	repository := NewFileRepository(fileSystem)
	readQuery, err := repository.GetMigrationQuery("/tmp/1.sql")

	if readQuery != "" || err == nil {
		test.Fail()
	}
}