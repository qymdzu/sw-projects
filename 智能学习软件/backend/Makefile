BINARY_NAME=smart-learning-server
GO=go

.PHONY: all build run test vet tidy clean

all: vet build test

build:
	$(GO) build -o bin/$(BINARY_NAME) ./cmd/server

run:
	$(GO) run ./cmd/server

test:
	$(GO) test ./... -count=1 -race

cover:
	$(GO) test ./... -count=1 -coverprofile=coverage.out -covermode=atomic
	$(GO) tool cover -html=coverage.out -o coverage.html

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

clean:
	rm -rf bin/ coverage.out coverage.html