MOD_NAME = $(shell go list -m)
BIN_NAME = $(shell basename $(MOD_NAME))

.PHONY: build
build: lint
	go build -race -o $(BIN_NAME)

.PHONY: lint
lint:
	go vet

.PHONY: clean
clean:
	go clean
	go mod tidy
	rm -rf out
