<<<<<<< HEAD
.PHONY: build run clean test deps build-all build-linux build-windows build-mac

# Variables
BINARY_NAME=webscraper
BUILD_DIR=build
VERSION=1.0

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

# Build for all major platforms
build-all: build-linux build-windows build-mac

# Build for Linux (64-bit)
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go

# Build for Windows (64-bit)  
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go

# Build for macOS (64-bit)
build-mac:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go

# Build for macOS (Apple Silicon)
build-mac-arm:
	@echo "Building for macOS (Apple Silicon)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go

# Create portable package
package: build-all
	@echo "Creating portable packages..."
	@mkdir -p $(BUILD_DIR)/packages
	
	# Linux package
	@mkdir -p $(BUILD_DIR)/packages/linux
	@cp $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(BUILD_DIR)/packages/linux/$(BINARY_NAME)
	@cp config.yaml $(BUILD_DIR)/packages/linux/
	@cp -r web $(BUILD_DIR)/packages/linux/ 2>/dev/null || mkdir -p $(BUILD_DIR)/packages/linux/web/templates
	@echo "#!/bin/bash" > $(BUILD_DIR)/packages/linux/run.sh
	@echo "mkdir -p data" >> $(BUILD_DIR)/packages/linux/run.sh
	@echo "./$(BINARY_NAME)" >> $(BUILD_DIR)/packages/linux/run.sh
	@chmod +x $(BUILD_DIR)/packages/linux/run.sh
	@cd $(BUILD_DIR)/packages && tar -czf $(BINARY_NAME)-linux-amd64.tar.gz linux/
	
	# Windows package
	@mkdir -p $(BUILD_DIR)/packages/windows
	@cp $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(BUILD_DIR)/packages/windows/$(BINARY_NAME).exe
	@cp config.yaml $(BUILD_DIR)/packages/windows/
	@cp -r web $(BUILD_DIR)/packages/windows/ 2>/dev/null || mkdir -p $(BUILD_DIR)/packages/windows/web/templates
	@echo "@echo off" > $(BUILD_DIR)/packages/windows/run.bat
	@echo "if not exist data mkdir data" >> $(BUILD_DIR)/packages/windows/run.bat
	@echo "$(BINARY_NAME).exe" >> $(BUILD_DIR)/packages/windows/run.bat
	@echo "pause" >> $(BUILD_DIR)/packages/windows/run.bat
	@cd $(BUILD_DIR)/packages && zip -r $(BINARY_NAME)-windows-amd64.zip windows/
	
	# macOS package
	@mkdir -p $(BUILD_DIR)/packages/macos
	@cp $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(BUILD_DIR)/packages/macos/$(BINARY_NAME)
	@cp config.yaml $(BUILD_DIR)/packages/macos/
	@cp -r web $(BUILD_DIR)/packages/macos/ 2>/dev/null || mkdir -p $(BUILD_DIR)/packages/macos/web/templates
	@echo "#!/bin/bash" > $(BUILD_DIR)/packages/macos/run.sh
	@echo "mkdir -p data" >> $(BUILD_DIR)/packages/macos/run.sh
	@echo "./$(BINARY_NAME)" >> $(BUILD_DIR)/packages/macos/run.sh
	@chmod +x $(BUILD_DIR)/packages/macos/run.sh
	@cd $(BUILD_DIR)/packages && tar -czf $(BINARY_NAME)-darwin-amd64.tar.gz macos/
	
	@echo "Portable packages created in $(BUILD_DIR)/packages/"

# Run the application
run:
	@echo "Starting $(BINARY_NAME)..."
	@mkdir -p data web/templates
	@go run main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f data/*.db

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Development setup
setup: deps
	@echo "Setting up development environment..."
	@mkdir -p data web/templates
=======
.PHONY: build run clean test deps build-all build-linux build-windows build-mac

# Variables
BINARY_NAME=webscraper
BUILD_DIR=build
VERSION=1.0

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

# Build for all major platforms
build-all: build-linux build-windows build-mac

# Build for Linux (64-bit)
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go

# Build for Windows (64-bit)  
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go

# Build for macOS (64-bit)
build-mac:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go

# Build for macOS (Apple Silicon)
build-mac-arm:
	@echo "Building for macOS (Apple Silicon)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go

# Create portable package
package: build-all
	@echo "Creating portable packages..."
	@mkdir -p $(BUILD_DIR)/packages
	
	# Linux package
	@mkdir -p $(BUILD_DIR)/packages/linux
	@cp $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(BUILD_DIR)/packages/linux/$(BINARY_NAME)
	@cp config.yaml $(BUILD_DIR)/packages/linux/
	@cp -r web $(BUILD_DIR)/packages/linux/ 2>/dev/null || mkdir -p $(BUILD_DIR)/packages/linux/web/templates
	@echo "#!/bin/bash" > $(BUILD_DIR)/packages/linux/run.sh
	@echo "mkdir -p data" >> $(BUILD_DIR)/packages/linux/run.sh
	@echo "./$(BINARY_NAME)" >> $(BUILD_DIR)/packages/linux/run.sh
	@chmod +x $(BUILD_DIR)/packages/linux/run.sh
	@cd $(BUILD_DIR)/packages && tar -czf $(BINARY_NAME)-linux-amd64.tar.gz linux/
	
	# Windows package
	@mkdir -p $(BUILD_DIR)/packages/windows
	@cp $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(BUILD_DIR)/packages/windows/$(BINARY_NAME).exe
	@cp config.yaml $(BUILD_DIR)/packages/windows/
	@cp -r web $(BUILD_DIR)/packages/windows/ 2>/dev/null || mkdir -p $(BUILD_DIR)/packages/windows/web/templates
	@echo "@echo off" > $(BUILD_DIR)/packages/windows/run.bat
	@echo "if not exist data mkdir data" >> $(BUILD_DIR)/packages/windows/run.bat
	@echo "$(BINARY_NAME).exe" >> $(BUILD_DIR)/packages/windows/run.bat
	@echo "pause" >> $(BUILD_DIR)/packages/windows/run.bat
	@cd $(BUILD_DIR)/packages && zip -r $(BINARY_NAME)-windows-amd64.zip windows/
	
	# macOS package
	@mkdir -p $(BUILD_DIR)/packages/macos
	@cp $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(BUILD_DIR)/packages/macos/$(BINARY_NAME)
	@cp config.yaml $(BUILD_DIR)/packages/macos/
	@cp -r web $(BUILD_DIR)/packages/macos/ 2>/dev/null || mkdir -p $(BUILD_DIR)/packages/macos/web/templates
	@echo "#!/bin/bash" > $(BUILD_DIR)/packages/macos/run.sh
	@echo "mkdir -p data" >> $(BUILD_DIR)/packages/macos/run.sh
	@echo "./$(BINARY_NAME)" >> $(BUILD_DIR)/packages/macos/run.sh
	@chmod +x $(BUILD_DIR)/packages/macos/run.sh
	@cd $(BUILD_DIR)/packages && tar -czf $(BINARY_NAME)-darwin-amd64.tar.gz macos/
	
	@echo "Portable packages created in $(BUILD_DIR)/packages/"

# Run the application
run:
	@echo "Starting $(BINARY_NAME)..."
	@mkdir -p data web/templates
	@go run main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f data/*.db

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Development setup
setup: deps
	@echo "Setting up development environment..."
	@mkdir -p data web/templates
>>>>>>> master
	@echo "Setup complete! Run 'make run' to start the application."