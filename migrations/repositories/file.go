package repositories

import (
	"fmt"
	"os"

	"github.com/jimenezmaximiliano/migrations/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/migrations/helpers"
)

// FileRepository fetches migrations files from a given path.
type FileRepository interface {
	GetMigrationFilePaths(migrationsDirectoryAbsolutePath string) ([]string, error)
	GetMigrationQuery(migrationAbsolutePath string) (string, error)
}

type fileRepository struct {
	fileSystem adapters.FileSystem
}

// Ensure fileRepository implements FileRepository
var _ FileRepository = fileRepository{}

// NewFileRepository returns an implementation of FileRepository.
func NewFileRepository(fileSystem adapters.FileSystem) FileRepository {
	return fileRepository{
		fileSystem: fileSystem,
	}
}

// GetMigrationFilePaths returns the paths of the migrations files on the given directory.
func (repository fileRepository) GetMigrationFilePaths(migrationsDirectoryAbsolutePath string) ([]string, error) {
	migrationsDirectoryAbsolutePath = helpers.AddTrailingSlashToPathIfNeeded(migrationsDirectoryAbsolutePath)
	migrationFiles, err := repository.fileSystem.ReadDir(migrationsDirectoryAbsolutePath)
	if err != nil {
		return nil, fmt.Errorf("could not read files from the migrations directory (path: %s)\n%w", migrationsDirectoryAbsolutePath, err)
	}

	return getMigrationFilePathsFromFiles(migrationFiles, migrationsDirectoryAbsolutePath), nil
}

// GetMigrationQuery returns the query for a migration file path.
func (repository fileRepository) GetMigrationQuery(migrationAbsolutePath string) (string, error) {
	query, err := repository.fileSystem.ReadFile(migrationAbsolutePath)
	if err != nil {
		return "", fmt.Errorf("could not read contents of a migration file (path: %s) \n%w", migrationAbsolutePath, err)
	}

	return string(query), nil
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
