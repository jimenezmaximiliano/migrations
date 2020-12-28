package migration

type MigrationCollection struct {
	migrations map[string]Migration
}

func (collection *MigrationCollection) Add(migration Migration) {
	if collection.migrations == nil {
		collection.migrations = make(map[string]Migration)
	}
	collection.migrations[migration.GetAbsolutePath()] = migration
}

func (collection *MigrationCollection) ContainsMigrationPath(migrationPath string) bool {
	if _, migration := collection.migrations[migrationPath]; migration {
		return true
	}

	return false
}

func (collection *MigrationCollection) GetAll() []Migration {

	migrations := []Migration{}

	for _, migration := range collection.migrations {
		migrations = append(migrations, migration)
	}

	return migrations
}

func (collection *MigrationCollection) IsEmpty() bool {
	return len(collection.migrations) == 0
}

func (collection *MigrationCollection) GetMigrationsToRun() []Migration {
	migrations := []Migration{}

	for _, migration := range collection.migrations {

		if !migration.ShouldBeRun() {
			continue
		}

		migrations = append(migrations, migration)
	}

	return migrations
}
