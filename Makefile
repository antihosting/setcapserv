APP_NAME := trd

BIN := $(APP_NAME)
BIN_MAC := $(BIN)_darwin
BIN_LINUX := $(BIN)_linux
BIN_WINDOWS := $(BIN).exe
TARGET := target

VERSION := $(shell git describe --tags --always --dirty)
NOW := $(shell date +"%m-%d-%Y")

build:
	go build -v -ldflags "-X main.Version=$(VERSION) -X main.Build=$(NOW)"

target: build
	rm -rf $(TARGET)
	mkdir $(TARGET)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_LINUX) -v -ldflags "-X main.Version=$(VERSION) -X main.Build=$(NOW)"
	mv $(BIN_LINUX) $(TARGET)/$(BIN_LINUX)
	gzip $(TARGET)/$(BIN_LINUX)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BIN_MAC) -v -ldflags "-X main.Version=$(VERSION) -X main.Build=$(NOW)"
	mv $(BIN_MAC) $(TARGET)/$(BIN_MAC)
	gzip $(TARGET)/$(BIN_MAC)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -o $(BIN_WINDOWS) -v -ldflags "-X main.Version=$(VERSION) -X main.Build=$(NOW)"
	mv $(BIN_WINDOWS) $(TARGET)/$(BIN_WINDOWS)
	gzip $(TARGET)/$(BIN_WINDOWS)

update:
	go get -u ./...

test:
	go test -v

clean:
	go clean
	rm -f $(BIN)
	rm -f $(BIN_LINUX)
	rm -f $(BIN_WINDOWS)
