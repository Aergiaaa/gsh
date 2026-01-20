.PHONY: all build clean test install uninstall run fmt vet

BINARY_NAME=gsh
GOFLAGS=
GOBIN=$(go env GOPATH)/bin

all: install

build: install

uninstall:
	rm -f $(GOBIN)/$(BINARY_NAME)

install:
	go build $(GOFLAGS) -o $(GOBIN)/$(BINARY_NAME)

clean:
	go clean
	rm -f $(BINARY_NAME)

test:
	go test -v ./...

run:
	go run .
