package db

import (
	"database/sql"
	"fmt"
	"time"
	"tt-copier/internal/fileutils"
	"tt-copier/internal/logger"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func NewDBInstance(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS uploaded_logs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            timestamp TEXT,
            source_path TEXT,
            dest_path TEXT,
            file_name TEXT
        );`

	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, fmt.Errorf("error creating uploaded_logs table: %v", err)
	}

	return &DB{db: db}, nil
}

func (l *DB) Close() error {
	return l.db.Close()
}

func (l *DB) LogEntry(sourcePath string, destinationPath string, fileName string) error {
	stmt, err := l.db.Prepare("INSERT INTO uploaded_logs (timestamp, source_path, dest_path, file_name) VALUES (?, ?, ?, ?)")

	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now().Format("01 Jan 2000 12:00:00"), sourcePath, destinationPath, fileName)

	if err != nil {
		return fmt.Errorf("error executing statement: %v", err)
	}

	return nil
}

func (l *DB) IsFileUploaded(fileName string) (bool, error) {
	var count int

	err := l.db.QueryRow("SELECT COUNT(*) FROM uploaded_logs WHERE file_name = ?", fileName).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("error querying file existence: %v", err)
	}

	logger.Info(fmt.Sprintf("File %s count is %d", fileName, count), "DB", "INFO")

	return count > 0, nil
}

func FilterUploadedFiles(dbInstance *DB, files []fileutils.LocalFileInfo) ([]fileutils.LocalFileInfo, error) {
	var filteredFiles []fileutils.LocalFileInfo

	for _, file := range files {
		isUploaded, err := dbInstance.IsFileUploaded(file.Name())

		if err != nil {
			return nil, err
		}
		if !isUploaded {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles, nil
}
