package repositories

import (
	"io/fs"
	"os"

	"github.com/pkg/errors"

	"github.com/jimenezmaximiliano/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/helpers"
)

// FileRepository fetches migrations files from a given path.
type FileRepository interface {
	GetMigrationFilePaths(migrationsDirectoryAbsolutePath string) ([]string, error)
	GetMigrationQuery(migrationAbsolutePath string) (string, error)
	CreateMigration(migrationAbsolutePath, query string) error
}

type fileRepository struct {
	fileSystem adapters.FileSystem
}

// Ensure fileRepository implements FileRepository.
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
		return nil, errors.Wrapf(
			err,
			"could not read files from the migrations directory [%s]",
			migrationsDirectoryAbsolutePath,
		)
	}

	return getMigrationFilePathsFromFiles(migrationFiles, migrationsDirectoryAbsolutePath), nil
}

// GetMigrationQuery returns the query for a migration file path.
func (repository fileRepository) GetMigrationQuery(migrationAbsolutePath string) (string, error) {
	query, err := repository.fileSystem.ReadFile(migrationAbsolutePath)
	if err != nil {
		return "", errors.Wrapf(err, "could not read contents of a migration file [%s]", migrationAbsolutePath)
	}

	return string(query), nil
}

// CreateMigration creates a new file with th emigration content.
func (repository fileRepository) CreateMigration(migrationAbsolutePath, query string) error {
	err := repository.fileSystem.WriteFile(migrationAbsolutePath, []byte(query), fs.FileMode(0644))
	if err != nil {
		return errors.Wrapf(err, "failed to create migration file on [%s]", migrationAbsolutePath)
	}

	return nil
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
