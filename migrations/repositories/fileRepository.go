package repositories

import (
	"fmt"
	"os"
)

// DirectoryReader is an interface that handles reading files from directories
type DirectoryReader interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
}

type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

// FileRepository is an interface to get migrations files from a given path
type FileRepository interface {
	GetMigrationFilePaths(migrationsDirectoryAbsolutePath string) ([]string, error)
	GetMigrationQuery(migrationAbsolutePath string) (string, error)
}

type fileRepository struct {
	directoryReader DirectoryReader
	fileReader      FileReader
}

// NewFileRepository allows you to get an implementation of FileRepository
func NewFileRepository(directoryReader DirectoryReader, fileReader FileReader) FileRepository {
	return fileRepository{
		directoryReader: directoryReader,
		fileReader:      fileReader,
	}
}

func (repository fileRepository) GetMigrationFilePaths(migrationsDirectoryAbsolutePath string) ([]string, error) {
	migrationsDirectoryAbsolutePath = addTrailingSlashIfNeeded(migrationsDirectoryAbsolutePath)

	migrationFiles, err := repository.directoryReader.ReadDir(migrationsDirectoryAbsolutePath)

	if err != nil {
		return nil, fmt.Errorf("migrations.readMigrationsPath (path: %s) \n%w", migrationsDirectoryAbsolutePath, err)
	}

	return getMigrationFilePathsFromFiles(migrationFiles, migrationsDirectoryAbsolutePath), nil
}

func (repository fileRepository) GetMigrationQuery(migrationAbsolutePath string) (string, error) {

	query, err := repository.fileReader.ReadFile(migrationAbsolutePath)

	if err != nil {
		return "", fmt.Errorf("migrations.getMigrationQuery (path: %s) \n%w", migrationAbsolutePath, err)
	}

	return string(query), nil
}

func addTrailingSlashIfNeeded(path string) string {
	lastCharacterIndex := len(path) - 1
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

	return file.IsDir() || fileNameLength <= 4 || fileName[(fileNameLength-4):] != ".sql"
}
