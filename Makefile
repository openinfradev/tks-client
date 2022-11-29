all: clean lint fmt build

clean:
	rm -rf ./bin

lint:
	golangci-lint run

fmt:
	go list -f '{{.Dir}}' ./... | xargs -L1 gofmt -w

build:
	GOFLAGS=-mod=mod go build -o bin/tks main.go
