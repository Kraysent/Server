all: build test

build:
	go build .

test: 
	go test -p 1 ./...

run:
	go run .

run-migrations: 
	cd postgre; pgmigrate -v -t 1 migrate -c 'dbname=testserverdb'