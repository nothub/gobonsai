MOD_NAME = $(shell go list -m)
BIN_NAME = $(shell basename $(MOD_NAME))
GIT_TAG  = $(shell git describe --tags --abbrev=0 --dirty --match v[0-9]* 2> /dev/null || echo "v0.0.0-indev")
VERSION  = $(GIT_TAG:v%=%)
LDFLAGS  = -ldflags="-X '$(MOD_NAME)/version=$(VERSION)'"

out/$(BIN_NAME): $(shell ls go.mod go.sum *.go)
	$(info dev build $(MOD_NAME) $(BIN_NAME) $(GIT_TAG) $(VERSION))
	go build $(LDFLAGS) -race -o out/$(BIN_NAME)

.PHONY: release
release: clean
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o out/$(BIN_NAME)-linux
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o out/$(BIN_NAME)-linux-arm64
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o out/$(BIN_NAME)-darwin
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o out/$(BIN_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o out/$(BIN_NAME)-windows.exe
	./debs.sh $(VERSION)

.PHONY: clean
clean:
	go clean
	go mod tidy
	rm -rf out

.PHONY: check
check:
	go vet
	go test -v -parallel $(shell grep -c -E "^processor.*[0-9]+" "/proc/cpuinfo") $(MOD_NAME)/...
