GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=modify-instance-charge-type
VERSION=0.1.0
ARCH=amd64
BINARY_LINUX=$(BINARY_NAME)$(VERSION).linux-$(ARCH)
BINARY_DARWIN=$(BINARY_NAME)$(VERSION).darwin-$(ARCH)
BINARY_WINDOWS=$(BINARY_NAME)$(VERSION).windows-$(ARCH).exe

all: test clean build
build: build-linux build-darwin build-windows

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_LINUX)
	rm -f $(BINARY_DARWIN)
	rm -f $(BINARY_WINDOWS)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v $(BINARY_NAME).go

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DARWIN) -v $(BINARY_NAME).go

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v $(BINARY_NAME).go
