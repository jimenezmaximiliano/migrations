package models

import (
	"sort"

	"github.com/pkg/errors"
)

// Collection is a set of implementations of Migration indexed by migration absolute path.
type Collection struct {
	migrations map[string]Migration
}

// Add adds a new MigrationContainer to the collection.
func (collection *Collection) Add(migration Migration) error {
	if collection.migrations == nil {
		collection.migrations = make(map[string]Migration)
	}

	for _, currentMigration := range collection.migrations {
		if currentMigration.GetOrder() == migration.GetOrder() {
			return errors.Errorf("two migrations cannot have the same order [%s] [%s]",
				currentMigration.GetName(),
				migration.GetName())
		}
	}

	collection.migrations[migration.GetAbsolutePath()] = migration

	return nil
}

// ContainsMigrationPath check if a given path is already in the collection.
func (collection *Collection) ContainsMigrationPath(migrationPath string) bool {
	if _, migration := collection.migrations[migrationPath]; migration {
		return true
	}

	return false
}

// GetAll returns all the migrations in the collection.
func (collection *Collection) GetAll() []Migration {
	migrations := []Migration{}
	for _, migration := range collection.migrations {
		migrations = append(migrations, migration)
	}
	sortMigrations(migrations)

	return migrations
}

// IsEmpty returns true if the collection is empty.
func (collection *Collection) IsEmpty() bool {
	return len(collection.migrations) == 0
}

// GetMigrationsToRun returns a list of migrations that has not been run yet.
func (collection *Collection) GetMigrationsToRun() []Migration {
	migrations := []Migration{}
	for _, migration := range collection.migrations {
		if !migration.ShouldBeRun() {
			continue
		}
		migrations = append(migrations, migration)
	}
	sortMigrations(migrations)

	return migrations
}

func sortMigrations(migrations []Migration) {
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ShouldBeRunFirst(migrations[j])
	})
}
