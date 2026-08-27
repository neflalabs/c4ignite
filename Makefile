.PHONY: all build test clean run

BINARY=bin/c4ignite

all: test build

build:
	docker run --rm -v "$$(pwd)":/app -w /app golang:1.23-alpine go build -ldflags="-s -w" -o $(BINARY) ./cmd/c4ignite

test:
	docker run --rm -v "$$(pwd)":/app -w /app golang:1.23-alpine go test -v ./...

cross-build:
	mkdir -p dist
	docker run --rm -v "$$(pwd)":/app -w /app golang:1.23-alpine sh -c '\
		GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/c4ignite-linux-amd64 ./cmd/c4ignite && \
		GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/c4ignite-linux-arm64 ./cmd/c4ignite && \
		GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/c4ignite-darwin-amd64 ./cmd/c4ignite && \
		GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/c4ignite-darwin-arm64 ./cmd/c4ignite && \
		GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/c4ignite-windows-amd64.exe ./cmd/c4ignite \
	'

clean:
	rm -rf bin/ dist/
