# Go SFTP File Transfer Project

This project is a Go-based application responsible for transferring files using SFTP. It features a structured approach with dedicated modules for configuration, database interaction, file utilities, and SFTP operations.

--- Project Status: [Under Development]

## Installation and Setup Instructions

### Requirements

- go 1.16.5
- sqlite3
- gcc

### Installaing Go

Ubuntu

````bash
    sudo apt install golang-go
    ```

Arch
```bash
    sudo pacman -S go
    ```

Red Hat, Fedora, CentOS
```bash
    sudo yum install go
    ```

### Project Setup

- 1. Clone the project
```bash
    git clone github.com/hamooo-d/sftp-file-transfer
    cd sftp-file-transfer
    ```

- 2. Enable CGO and install GCC
```bash
    export CGO_ENABLED=1
    sudo apt install build-essential # For Ubuntu/Debian
    sudo yum groupinstall "Development Tools" # For RedHat/CentOS
    sudo pacman -S base-devel # For Arch Linux
    ```

- 3. Install dependencies
```bash
    go mod tidy
    ```

### Folder Structure
```bash
        /
        ├── cmd
        │   └── copier.go
        ├── config
        │   └── config.go
        ├── internal
        │   ├── db
        │   ├── fileutils
        │   ├── logger
        │   └── sftp
        ├── scripts
        │   ├── gen_mock_src_files.bash
        │   ├── redo_sftp.bash
        │   └── sftp_setup.bash
        ├── config.yaml
        ├── copier
        ├── db.sqlite
        ├── go.mod
        ├── go.sum
        ├── logs.log
        └── Makefile
````

The project is organized into several directories, each serving a specific purpose in the application's architecture:

- `cmd`: Contains the main application code.

  - `copier.go`: The main Go file for the application.

- `config`: Holds configuration related code.

  - `config.go`: Go file for handling configuration logic.

- `internal`: Consists of various internal packages used by the application.

  - `db`: Contains code related to database operations.
  - `fileutils`: Includes utilities for file handling.
  - `logger`: Manages logging functionalities.
  - `sftp`: Contains code for SFTP operations.

- `scripts`: Stores various scripts for setup and utility purposes durning the development stage.

  - `gen_mock_src_files.bash`: Generates mock source files for development.
  - `redo_sftp.bash`: Script to redo SFTP configurations.
  - `sftp_setup.bash`: Setup a mock SFTP folder and users/group structure for development.

- Root Directory Files:
  - `config.yaml`: The configuration file in YAML format.
  - `copier`: The compiled binary of the application.
  - `db.sqlite`: SQLite database file.
  - `go.mod` and `go.sum`: Go module files for managing dependencies.
  - `logs.log`: Log file for the application.
  - `Makefile`: Defines commands for building, setting up, and running the application.

### Makefile

- To build the project

````bash
    make build
    ```
- To build and run the project
```bash
    make run
    ```
- To run sftp server mock structure
```bash
    make setup
    ```
- To redo sftp server mock structure
```bash
    make redo
    ```
- To generate mock files
```bash
    make gen_mock
    ```

## Configuration

- Place the configuration file in the root directory of the project with the name `config.yaml`
```yaml
    # config.yaml
    db:
        path: "db.sqlite"
        table_name: "files"
    log:
        path: "logs.log"
````
