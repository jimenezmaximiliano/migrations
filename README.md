# Migrations

![migrate](https://bestanimations.com/media/birds/1460382957ducks-flying-gif.gif)

Migrations is a database migration tool that uses go's **database/sql** from the standard library

## How to use it

    go get github.com/jimenezmaximiliano/migrations

```golang
package main

import (
	"database/sql"

	// This example uses mysql but you can pick any other driver compatible with database/sql
	_ "github.com/go-sql-driver/mysql"

	"github.com/jimenezmaximiliano/migrations"
)

func main() {
	migrations.RunMigrationsCommand(func() (*sql.DB, error) {
		// Here you can set up your db connection
		return sql.Open("mysql", "user:password@/db?multiStatements=true")
	})
}
```

Then use the binary like this:

    ./migrate -path=/absolute/path/to/migrations

> See the [example directory](https://github.com/jimenezmaximiliano/migrations/tree/master/example) in this repository for a working example

### Migration files

All migration files must:

- be in the provided path (not inside subdirectories)
- end in *.sql* (files without the .sql extension will be ignored)
- be in this format: {number}_{string}.sql where number determines the order on which migrations will be run
- be ordered by filename in the order they should run 
  (you can use migo for that)

Example:

```bash
/app/migrations
/app/migrations/1627676712447528000_createGophersTable.sql
/app/migrations/1627676757857350000_createGolfersTable.sql
```

> See [example migrations](https://github.com/jimenezmaximiliano/migrations/tree/master/example/migrations) in the example directory

### Using migo to create migration files

Migo creates migration files using the current timestamp as a prefix.

#### Installation

```bash
go get -u github.com/jimenezmaximiliano/migrations/migo
```

#### Usage

```bash
migo -path=/app/migrations name=myMigration migration:create
```

The above command results in:

```bash
/app/migrations/1627676757857350000_myMigration.sql
```

## How it works / features

- migration files can contain multiple queries
- each successful migration will be added to the migrations table
- the migrations table is created automatically when migrations are run

## Customization

You can use the [migrations facade](https://github.com/jimenezmaximiliano/migrations/blob/master/facade.go)
as a tutorial on how to replace any component of the package by implementing one of its
interfaces.

## Breaking changes

To avoid go modules issues with major version, this package doesn't support semver.
Breaking changes could happen in minor versions but not in patches.