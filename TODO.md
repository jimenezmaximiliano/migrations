# TODO

## 1.0.0

- RunMigrations should return general errors on the return value and migration specific errors should be in the migration itself
- Support a trailing slash at the end of the path flag
- Clean up the code
- Write tests
- Setup github workflows
- Write docs
- Check that migration file names are unique
- Check where I'm adding context to errors
- Fail early if we cannot connect to the db
- Fail early if the migrations path is invalid
- Command status code
- Include binaries with support for several db drivers

## Future versions

- Add a verbose mode
- Allow reverting migrations