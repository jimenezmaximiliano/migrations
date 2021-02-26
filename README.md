# Migrations

Run example ./migrate -path=/home/maxi/projects/very-simple-migrations/example/migrations

./mockery --dir=migrations/adapters --name=DB
./mockery --dir=migrations/adapters --name=DBRows
./mockery --dir=migrations/adapters --name=FileSystem
./mockery --dir=migrations/adapters --name=File
./mockery --dir=migrations/adapters --name=OptionParser