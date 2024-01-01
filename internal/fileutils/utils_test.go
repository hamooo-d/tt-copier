package fileutils

import (
	"os"
	"strings"
	"testing"
	"time"
)

type mockFileInfo struct {
	modTime time.Time
	name    string
	size    int64
	mode    os.FileMode
	isDir   bool
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return nil }

func TestAddBankDestination(t *testing.T) {
	files := []os.FileInfo{
		mockFileInfo{name: "000003_sample.txt"},
		mockFileInfo{name: "000004_sample.txt"},
		mockFileInfo{name: "non_bank_file.txt"},
	}
	basePath := "/home/sftp/files/TTP"

	result, err := AddBankDestination(files, basePath)
	if err != nil {
		t.Fatalf("AddBankDestination returned an error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 files to match bank IDs, got %d", len(result))
		for _, f := range result {
			t.Logf("File: %s, Destination: %s", f.Name(), f.DestinationFullPath)
		}
	}

	for _, f := range result {
		if !strings.HasPrefix(f.DestinationFullPath, basePath) {
			t.Errorf("File destination path incorrect: %s", f.DestinationFullPath)
		}
	}

	for _, f := range result {
		t.Logf("File: %s, Destination: %s", f.Name(), f.DestinationFullPath)
	}

}

func TestAddGetDestination(t *testing.T) {
	files := []os.FileInfo{
		mockFileInfo{name: "sample1.txt"},
		mockFileInfo{name: "sample2.txt"},
	}
	destination := "/online/mxpprod/selectsystem_files/cardholder/in"

	result, err := AddGetDestination(files, destination)
	if err != nil {
		t.Errorf("AddGetDestination returned an error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 files, got %d", len(result))
	}
}
