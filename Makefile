MOD_NAME = $(shell go list -m)
BIN_NAME = $(shell basename $(MOD_NAME))

out/$(BIN_NAME): $(shell ls go.mod go.sum *.go)
	go build -race -o out/$(BIN_NAME)

.PHONY: lint
lint:
	go vet

.PHONY: check
check:
	go test -v -parallel $(shell grep -c -E "^processor.*[0-9]+" "/proc/cpuinfo") $(MOD_NAME)/...
