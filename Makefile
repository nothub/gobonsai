MOD_NAME = $(shell go list -m)
BIN_NAME = $(shell basename $(MOD_NAME))

out/$(BIN_NAME): $(shell ls go.mod go.sum *.go)
	go build -race -o out/$(BIN_NAME)

.PHONY: release
release: clean
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o out/$(BIN_NAME)-linux
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o out/$(BIN_NAME)-linux-arm64
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o out/$(BIN_NAME)-darwin
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o out/$(BIN_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o out/$(BIN_NAME)-windows.exe

.PHONY: clean
clean:
	go clean
	go mod tidy
	rm -rf out
