all: build test

build:
	go build .

test: 
	go test -p 1 ./...

run:
	go run .

run-migrations: 
	migrate -path postgre/migrations -database "postgres://localhost:5432/testserverdb?sslmode=disable&user=testserver&password=passw0rd" -verbose up