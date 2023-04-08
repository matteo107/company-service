# Go binary name
EXECUTABLE=companysrv
VERSION=$(shell git describe --tags --always --long --dirty)
WINDOWS=$(EXECUTABLE).exe
LINUX=$(EXECUTABLE)
SRC_FILES=cmd/api
GO_COMPILER=go
GO_FLAGS=-v

.PHONY: all build clean

test:
	go test -v ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

windows:
	go build -v -o bin/$(WINDOWS)  $(SRC_FILES)

linux:
	go build -v -o bin/$(LINUX)  $(SRC_FILES)

# Clean the artifacts
clean:
	rm -f bin/$(EXECUTABLE)

# Build and install the binary
all: clean windows linux install

# Build the binary
build: windows linux

# Install the binary
install:
	go install $(GO_FLAGS) $(SRC_FILES)

# Run the binary
run:
	docker-compose up --build

# Show the version
version:
	@echo $(VERSION)