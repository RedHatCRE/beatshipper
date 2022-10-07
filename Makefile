all: test build

test:
	go test ./...

build:
	go build -v -o builds/gz-beat-shipper cmd/shipper/main.go