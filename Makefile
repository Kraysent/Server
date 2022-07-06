all: build test

build:
	go build .

test: 
	go test -p 1 ./...

run:
	go run .
