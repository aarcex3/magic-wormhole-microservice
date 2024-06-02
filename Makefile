# Variables
SRC_DIR=./cmd
SRC=$(SRC_DIR)/main.go
BINARY=$(SRC_DIR)/main

# Default target: build and run the application
all: run

# Build the binary
build:
	go build -o $(BINARY) $(SRC)

# Run the application
run: build
	./$(BINARY)

# Clean up binaries
clean:
	rm -f $(BINARY)

# Phony targets
.PHONY: all build run clean
