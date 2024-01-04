SETUP_SCRIPT = ./scripts/sftp_setup.sh
REDO_SCRIPT = ./scripts/sftp_redo.sh
GO_FILE = ./cmd/copier.go
BINARY_NAME = copier

all: build

setup:
	@echo "Running setup script..."
	@bash $(SETUP_SCRIPT)

redo:
	@echo "Running redo script..."
	@bash $(REDO_SCRIPT)
	
gen_mock:
	@echo "Generating mocks..."
	@./scripts/gen_mock_src_files.bash

build:
	@echo "Building Go file..."
	@go build -o $(BINARY_NAME) $(GO_FILE)

run: build
	@echo "Running Go application..."
	@./$(BINARY_NAME)
