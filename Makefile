BINARY_NAME=cli-convert

GO_BUILD=go build
GO_TEST=go test

.PHONY: all build test clean

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO_BUILD) -o $(BINARY_NAME) .
	@echo "$(BINARY_NAME) built successfully."

test:
	@echo "Running tests..."
	$(GO_TEST) -v ./...

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	@echo "Cleaned successfully."
