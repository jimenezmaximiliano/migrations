package fileRepository

import (
	"fmt"
	"os"
)

type DirectoryReader interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
}

type FileRepository interface {
	GetMigrationFilePaths(migrationsDirectoryAbsolutePath string) ([]string, error)
}

type fileRepository struct {
	directoryReader DirectoryReader
}

func New(reader DirectoryReader) FileRepository {
	return fileRepository{
		directoryReader: reader,
	}
}

func (repository fileRepository) GetMigrationFilePaths(migrationsDirectoryAbsolutePath string) ([]string, error) {
	migrationsDirectoryAbsolutePath = addTrailingSlashIfNeeded(migrationsDirectoryAbsolutePath)

	migrationFiles, err := repository.directoryReader.ReadDir(migrationsDirectoryAbsolutePath)

	if err != nil {
		return []string{}, fmt.Errorf("verySimpleMigrations.readMigrationsPath (path: %s) \n%w", migrationsDirectoryAbsolutePath, err)
	}

	return getMigrationFilePathsFromFiles(migrationFiles, migrationsDirectoryAbsolutePath), nil
}

func addTrailingSlashIfNeeded(path string) string {
	lastCharacterIndex := len(path) -1
	if path[lastCharacterIndex:] != "/" {
		return path + "/"
	}

	return path
}

func getMigrationFilePathsFromFiles(files []os.FileInfo, migrationsDirectoryAbsolutePath string) []string {
	var migrationFilePaths []string

	for _, file := range files {
		if isNotASqlFile(file) {
			continue
		}

		currentMigrationAbsolutePath := migrationsDirectoryAbsolutePath + file.Name()

		migrationFilePaths = append(migrationFilePaths, currentMigrationAbsolutePath)
	}

	return migrationFilePaths
}

func isNotASqlFile(file os.FileInfo) bool {
	fileName := file.Name()
	fileNameLength := len(fileName)

	return file.IsDir() || fileNameLength <= 4 || fileName[(fileNameLength - 4):] != ".sql"
}