# Migrations

Run example ./migrate -path=/home/maxi/projects/very-simple-migrations/example/migrations

./mockery --dir=adapters --name=DB
./mockery --dir=adapters --name=DBRows
./mockery --dir=adapters --name=FileSystem
./mockery --dir=adapters --name=File
./mockery --dir=adapters --name=OptionParser
./mockery --dir=repositories --name=DBRepository
./mockery --dir=repositories --name=FileRepository
./mockery --dir=services --name=Fetcher