# Makefile
APP_NAME = sia
OUTPUT_DIR = bin

build: clean
	# Build for Windows (64-bit)
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(APP_NAME)-windows.exe
	# Build for macOS (Intel x86, 64-bit)
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(APP_NAME)-mac-intel
	# Build for macOS (ARM, 64-bit)
	GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(APP_NAME)-mac-arm
	# Build for Linux (64-bit)
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(APP_NAME)-linux
	# Build for Linux (ARM 64-bit)
	GOOS=linux GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(APP_NAME)-linux-arm
	# Build for Linux (ARMv7)
	GOOS=linux GOARCH=arm GOARM=7 go build -o $(OUTPUT_DIR)/$(APP_NAME)-linux-armv7

clean:
	rm -rf $(OUTPUT_DIR)
	mkdir -p $(OUTPUT_DIR)