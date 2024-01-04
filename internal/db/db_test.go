package db

import (
	"os"
	"testing"
	"time"
	"tt-copier/internal/fileutils"
)

func setupTestDB(t *testing.T) *DB {
	t.Helper()

	dbFile := "test_db.sqlite"
	db, err := NewDBInstance(dbFile)
	if err != nil {
		t.Fatalf("Failed to create test DB instance: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.Remove(dbFile)
	})

	return db
}

func TestNewDBInstance(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Errorf("Expected new DB instance, got nil")
	}
}

func TestLogEntry(t *testing.T) {
	db := setupTestDB(t)
	err := db.LogEntry("/source/path", "/dest/path", "testfile.txt", "COPY")
	if err != nil {
		t.Errorf("Failed to log entry: %v", err)
	}
}

func TestFileExists(t *testing.T) {
	db := setupTestDB(t)
	db.LogEntry("/source/path", "/dest/path", "existfile.txt", "COPY")

	exists, err := db.FileExists("existfile.txt")
	if err != nil {
		t.Errorf("Error checking file existence: %v", err)
	}
	if !exists {
		t.Errorf("File should exist, but it does not")
	}
}

func TestFileIsRenamed(t *testing.T) {
	db := setupTestDB(t)
	err := db.LogEntry("/source/path", "/dest/path", "renamedfile.txt", "Rename")
	if err != nil {
		t.Fatalf("Failed to log entry: %v", err)
	}

	renamed, err := db.FileIsRenamed("renamedfile.txt")
	if err != nil {
		t.Errorf("Error checking file rename status: %v", err)
	}
	if !renamed {
		t.Errorf("File should be marked as renamed, but it is not. Count should be > 0.")
	}
}

func TestFilterNotRenamed(t *testing.T) {
	db := setupTestDB(t)
	err := db.LogEntry("/source/path", "/dest/path", "renamedfile.txt", "Rename")
	if err != nil {
		t.Fatalf("Failed to log entry: %v", err)
	}

	files := []fileutils.LocalFileInfo{
		{FileInfo: mockFileInfo{name: "renamedfile.txt"}},
		{FileInfo: mockFileInfo{name: "notrenamed.txt"}},
	}

	filteredFiles, err := FilterNotRenamed(db, files)
	if err != nil {
		t.Fatalf("Error filtering not renamed files: %v", err)
	}

	if len(filteredFiles) != 1 || filteredFiles[0].Name() != "notrenamed.txt" {
		t.Errorf("FilterNotRenamed did not return the expected files")
		for _, file := range filteredFiles {
			t.Logf("Returned File: %s", file.Name())
		}
	}
}

func TestKeepNotPut(t *testing.T) {
	db := setupTestDB(t)
	db.LogEntry("/source/path", "/dest/path", "putfile.txt", "PUT")

	files := []fileutils.LocalFileInfo{
		{FileInfo: mockFileInfo{name: "putfile.txt"}},
		{FileInfo: mockFileInfo{name: "notput.txt"}},
	}

	filteredFiles, err := KeepNotPut(db, files)
	if err != nil {
		t.Fatalf("Error filtering not put files: %v", err)
	}

	if len(filteredFiles) != 1 || filteredFiles[0].Name() != "notput.txt" {
		t.Errorf("KeepNotPut did not return the expected files")
	}
}

type mockFileInfo struct {
	name string
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() os.FileMode  { return 0 }
func (m mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m mockFileInfo) IsDir() bool        { return false }
func (m mockFileInfo) Sys() interface{}   { return nil }
