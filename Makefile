# Makefile for dill-monitor

# Binary name
BINARY_NAME=dill-monitor

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Detect operating system and architecture
OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

# Determine home directory based on operating system
ifeq ($(OS),windows)
	HOME_DIR=$(USERPROFILE)
	INSTALL_DIR=$(HOME_DIR)\.dill_monitor
	MKDIR=mkdir -p
	BINARY_EXT=.exe
	CP=copy
	RM=del /f /q
	CONFIG_DIR=config
else
	HOME_DIR=$(HOME)
	INSTALL_DIR=$(HOME_DIR)/.dill_monitor
	MKDIR=mkdir -p
	BINARY_EXT=
	CP=cp -f
	RM=rm -f
	CONFIG_DIR=config
	ifeq ($(OS),darwin)
		# macOS specific settings if needed
	endif
endif

# Build targets for different OS
.PHONY: all build build-windows build-linux build-darwin install clean test

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME)$(BINARY_EXT) -v ./cmd/server

# Cross compilation targets
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME).exe -v ./cmd/server

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_linux -v ./cmd/server

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_darwin -v ./cmd/server

build-all: build-windows build-linux build-darwin

# Installation
install: build
	@echo "Installing dill-monitor..."
	$(MKDIR) $(INSTALL_DIR)
	@echo "Copying binary to $(INSTALL_DIR)..."
ifeq ($(OS),windows)
	$(CP) $(BINARY_NAME)$(BINARY_EXT) $(INSTALL_DIR)
	@echo "Copying configuration files..."
	-$(CP) $(CONFIG_DIR)\config.json $(INSTALL_DIR)
	-$(CP) $(CONFIG_DIR)\server_config.json $(INSTALL_DIR)
else
	$(CP) $(BINARY_NAME)$(BINARY_EXT) $(INSTALL_DIR)/
	@echo "Copying configuration files..."
	-$(CP) $(CONFIG_DIR)/config.json $(INSTALL_DIR)/
	-$(CP) $(CONFIG_DIR)/server_config.json $(INSTALL_DIR)/
	chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
endif
	@echo "Installation completed!"
	@echo "Configuration files are located at: $(INSTALL_DIR)"

# For system-wide installation (Linux/macOS only)
install-system: build
ifneq ($(OS),windows)
	sudo $(CP) $(BINARY_NAME) /usr/local/bin/
	sudo $(MKDIR) /etc/dill_monitor
	-sudo $(CP) $(CONFIG_DIR)/config.json /etc/dill_monitor/
	-sudo $(CP) $(CONFIG_DIR)/server_config.json /etc/dill_monitor/
	@echo "System-wide installation completed!"
	@echo "Configuration files are located at: /etc/dill_monitor/"
else
	@echo "System-wide installation is not supported on Windows"
endif

# Cleanup
clean:
	@echo "Cleaning up..."
ifeq ($(OS),windows)
	$(RM) $(BINARY_NAME).exe
	$(RM) $(BINARY_NAME)_linux
	$(RM) $(BINARY_NAME)_darwin
else
	$(RM) $(BINARY_NAME)
	$(RM) $(BINARY_NAME).exe
	$(RM) $(BINARY_NAME)_linux
	$(RM) $(BINARY_NAME)_darwin
endif
	$(GOCLEAN)
	@echo "Cleanup completed!"

# Run tests
test:
	$(GOTEST) -v ./...
