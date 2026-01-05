PROJECT="boc/backend"
BINARY_NAME=boc-backend

all: clean tidy build run

.PHONY: build run tidy test clean

run:
	go run ./cmd/server

build:
	CGO_ENABLED=0 go build \
    		--tags local \
    		-v \
    		-o build/$(BINARY_NAME) \
    		./cmd/server

tidy:
	go mod tidy

test:
	CGO_ENABLED=0 go test -v ./...

clean:
	go clean
	rm -f build/$(BINARY_NAME)