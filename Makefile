all: build test

build:
	go build .

test: 
	go test -p 1 ./...

run:
	go run .

run-migrations: 
	migrate -path postgre/migrations -database "$(DB_URL)" -verbose up