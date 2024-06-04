# Variables
SRC_DIR=./cmd
SRC=$(SRC_DIR)/main.go
BINARY=$(SRC_DIR)/main

# Default target: build and run the application
all: run

# Generate code using templ
generate:
	templ generate

# Build the binary
build: generate
	go build -o $(BINARY) $(SRC)

# Run the application
run: build
	./$(BINARY)

# Clean up binaries
clean:
	rm -f $(BINARY)

# Phony targets
.PHONY: all generate build run clean
