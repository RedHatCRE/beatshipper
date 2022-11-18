all: test build

test:
	go test ./...

build:
	go build -v -o builds/beatshipper main.go