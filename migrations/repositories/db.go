package repositories

import (
	"fmt"

	"github.com/jimenezmaximiliano/migrations/migrations/adapters"
)

// DBRepository runs migration queries and handles the migrations table.
type DBRepository interface {
	CreateMigrationsTableIfNeeded() error
	GetAlreadyRunMigrationFilePaths(migrationsDirectoryAbsolutePath string) ([]string, error)
	RunMigrationQuery(query string) error
	RegisterRunMigration(migrationFileName string) error
	Ping() error
}

type dbRepository struct {
	db adapters.DB
}

// Ensure dbRepository implements DBRepository
var _ DBRepository = dbRepository{}

// NewDbRepository returns an implementation of DbRepository.
func NewDbRepository(db adapters.DB) DBRepository {
	return dbRepository{
		db: db,
	}
}

// Ping the DB to check if the connection is working
func (repository dbRepository) Ping() error {
	err := repository.db.Ping()
	if err == nil {
		return nil
	}

	return fmt.Errorf("could not connect to the DB\n%w", err)
}

// CreateMigrationsTableIfNeeded creates the migrations table used to keep track of already run migrations.
func (repository dbRepository) CreateMigrationsTableIfNeeded() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTO_INCREMENT,
			migration TEXT
		);`
	_, err := repository.db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not create the migrations table: \n%w", err)
	}

	return nil
}

// GetAlreadyRunMigrationFilePaths returns a list of migration file paths that have been run already.
func (repository dbRepository) GetAlreadyRunMigrationFilePaths(migrationsDirectoryAbsolutePath string) ([]string, error) {
	rows, err := repository.db.Query("SELECT migration FROM migrations")
	if err != nil {
		return nil, fmt.Errorf("could not get already run migrations from the migrations table\n%w", err)
	}
	defer rows.Close()

	return getMigrationPathsFromRows(rows, migrationsDirectoryAbsolutePath)
}

// RunMigrationQuery runs the migration query.
func (repository dbRepository) RunMigrationQuery(query string) error {
	_, err := repository.db.Exec(query)

	return err
}

// RegisterRunMigration creates a record on the migrations table for a successfully run migration.
func (repository dbRepository) RegisterRunMigration(migrationFileName string) error {
	_, err := repository.db.Exec("INSERT INTO migrations (migration) VALUES (?)", migrationFileName)

	return err
}

func getMigrationPathsFromRows(rows adapters.DBRows, migrationsDirectoryAbsolutePath string) ([]string, error) {
	var migrationsAlreadyRun []string
	for rows.Next() {
		migrationFileName := ""
		err := rows.Scan(&migrationFileName)
		if err != nil {
			return migrationsAlreadyRun, fmt.Errorf("migrations.readMigrationRowFromMigrationsTable \n%w", err)
		}
		currentMigrationAbsolutePath := migrationsDirectoryAbsolutePath + migrationFileName
		migrationsAlreadyRun = append(migrationsAlreadyRun, currentMigrationAbsolutePath)
	}

	return migrationsAlreadyRun, nil
}
