all: build test

build:
	go build .

test: 
	go test -p 1 ./...

run:
	go run .

run-migrations: 
	cd postgre; pgmigrate -v -t 1 migrate -c 'host=localhost dbname=testserverdb user=testserver port=5432 password=passw0rd'