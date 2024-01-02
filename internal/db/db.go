package db

import (
	"database/sql"
	"fmt"
	"time"
	"tt-copier/internal/fileutils"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type DB struct {
	db *sql.DB
}

func NewDBInstance(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Optionally, here you can ensure the table 'logs' exists or create it if it does not.

	return &DB{db: db}, nil
}

// LogEntry logs a new entry in the database
func (l *DB) LogEntry(sourcePath, destinationPath, fileName, action string) error {
	stmt, err := l.db.Prepare("INSERT INTO logs (Timestamp, Source, Destination, Filename, Action) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now().Format("02 Jan 2006 15:04:05"), sourcePath, destinationPath, fileName, action)
	if err != nil {
		return fmt.Errorf("error executing statement: %v", err)
	}

	return nil
}

// FileExists checks if a file with the given filename exists in the log
func (l *DB) FileExists(fileName string) (bool, error) {
	var count int
	err := l.db.QueryRow("SELECT COUNT(*) FROM logs WHERE Filename = ?", fileName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error querying file existence: %v", err)
	}
	return count > 0, nil
}

// FileIsRenamed checks if a file with the given filename has a 'Rename' action in the log
func (l *DB) FileIsRenamed(fileName string) (bool, error) {
	var count int
	err := l.db.QueryRow("SELECT COUNT(*) FROM logs WHERE Filename = ? AND Action = 'Rename'", fileName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error querying file rename status: %v", err)
	}
	return count > 0, nil
}

// Close closes the database connection
func (l *DB) Close() error {
	return l.db.Close()
}
func FilterNotRenamed(dbInstance *DB, files []fileutils.LocalFileInfo) ([]fileutils.LocalFileInfo, error) {
	var filteredFiles []fileutils.LocalFileInfo
	for _, file := range files {
		renamed, err := dbInstance.FileIsRenamed(file.Name())
		if err != nil {
			return nil, err
		}
		if !renamed {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, nil
}

func KeepNotPut(dbInstance *DB, files []fileutils.LocalFileInfo) ([]fileutils.LocalFileInfo, error) {
	var filteredFiles []fileutils.LocalFileInfo
	for _, file := range files {
		put, err := dbInstance.FileExists(file.Name())
		if err != nil {
			return nil, err
		}
		if !put {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, nil
}
