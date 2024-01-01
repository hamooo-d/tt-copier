# Makefile

# Variables for script paths
SETUP_SCRIPT = /scripts/sftp_setup.sh
REDO_SCRIPT = /scripts/sftp_redo.sh

# Default target
all: setup

# Setup command
setup:
	@echo "Running setup script..."
	@bash $(SETUP_SCRIPT)

# Redo command
redo:
	@echo "Running redo script..."
	@bash $(REDO_SCRIPT)
