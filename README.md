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
- end in *.sql*
- be ordered by filename in the order they should run

Example:

```bash
/tmp/migrations
/tmp/migrations/20200318001000_createGophersTable.sql
/tmp/migrations/20200418001000_createGolfersTable.sql
```

> See [example migrations](https://github.com/jimenezmaximiliano/migrations/tree/master/example/migrations) in the example directory

## How it works / features

- migration files can contain multiple queries
- each successful migration will be added to the migrations table
- the migrations table is created automatically when migrations are run

## Customization

You can use the [migrations facade](https://github.com/jimenezmaximiliano/migrations/blob/master/facade.go)
as a tutorial on how to replace any component of the package by implementing one of its
interfaces.