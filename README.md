# Migrations

![migrate](https://bestanimations.com/media/birds/1460382957ducks-flying-gif.gif)

Migrations is a minimalistic database migration tool that uses go's **database/sql** from the standard library.

## Features

- supports [any database driver](https://github.com/golang/go/wiki/SQLDrivers) that is compatible with **database/sql**, including: MySQL, Microsoft SQL Server, PostgreSQL, Oracle and SQLite.
- migrations are simple SQL files
- migrations can contain multiple queries
- easy generation of migration files
- we keep track of run migrations
- minimal dependencies
- customizable
- support for environment variables

## Usage

### Examples

#### Create a migration file

```bash
./migrations create -name=createGophersTable -path=/app/migrations/
```

or use environment variables:

```bash
MIGRATIONS_COMMAND="create" MIGRATIONS_PATH="/app/migrations/" MIGRATIONS_NEW_MIGRATION_NAME="createGophersTable" ./migrations
```

#### Run migrations
```bash
./migrations migrate -path=/app/migrations/
```

or use environment variables:

```bash
MIGRATIONS_COMMAND="migrate" MIGRATIONS_PATH="/app/migrations/" ./migrations
```

### create command

The **create** command creates a file with a prefix using the current timestamp. That's going to be used to determine
the order on which migrations should be run.

For example, running:

```bash
./migrations create -name=createGophersTable -path=/app/migrations/
./migrations create -name=createGolfersTable -path=/app/migrations/
```

Will result in:

```bash
/app/migrations/1627676712447528000_createGophersTable.sql
/app/migrations/1627676757857350000_createGolfersTable.sql
```

### migrate command

The **migrate** command runs migrations and then displays a report with the result of each migration run, if any.

```bash
./migrations migrate -path=/app/migrations/
```

Example output:

```bash
[ INFO ] Run migrations
[  OK  ] 1627676712447528000_createGophersTable.sql
[  OK  ] 1627676757857350000_createGolfersTable.sql
[ INFO ] Done
```

## Setup

1) Get the module
```bash
go get github.com/jimenezmaximiliano/migrations
```

2) Create a go file using the next template:
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

3) Then use the binary like this:
```bash
./migrations migrate -path=/app/migrations
```
or
```bash
go run migrations.go migrate -path=/app/migrations
```

> See the [example directory](https://github.com/jimenezmaximiliano/migrations/tree/master/example) in this repository for a working example

### Migration files

All migration files must:

- be in the provided path (not inside subdirectories)
- end in *.sql* (files without the .sql extension will be ignored)
- be in this format: {number}_{string}.sql where number determines the order on which migrations will be run
- be ordered by filename in the order they should run 
  (you can use the **create** command for that)

Example:

```bash
/app/migrations
/app/migrations/1627676712447528000_createGophersTable.sql
/app/migrations/1627676757857350000_createGolfersTable.sql
```

> See [example migrations](https://github.com/jimenezmaximiliano/migrations/tree/master/example/migrations) in the example directory

## Customization

You can use the [migrations facade](https://github.com/jimenezmaximiliano/migrations/blob/master/facade.go)
as a tutorial on how to replace any component of the package by implementing one of its
interfaces.
