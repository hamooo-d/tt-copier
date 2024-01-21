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
	files := []LocalFileInfo{
		{FileInfo: mockFileInfo{name: "000006_EV_MERC_20240102.txt"}, Path: "", SourceFullPath: ""},
		{FileInfo: mockFileInfo{name: "000004_sample.txt"}, Path: "", SourceFullPath: ""},
		{FileInfo: mockFileInfo{name: "CL.000005.240123"}, Path: "", SourceFullPath: ""},
		{FileInfo: mockFileInfo{name: "non_bank_file.txt"}, Path: "", SourceFullPath: ""},
	}
	basePath := "/home/sftp/files/TTP"

	banksNames := map[string]string{
		"000001": "TT",
		"000002": "ATIB",
		"000003": "SB",
		"000004": "NAB",
		"000005": "MED",
		"000006": "NCB",
	}

	result, err := AddBankDestination(files, basePath, banksNames, "UAT")

	if err != nil {
		t.Fatalf("AddBankDestination returned an error: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 files to match bank IDs, got %d", len(result))

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
	files := []LocalFileInfo{
		{FileInfo: mockFileInfo{name: "prefix1_sample.txt"}, Path: "", SourceFullPath: ""},
		{FileInfo: mockFileInfo{name: "prefix2_sample.txt"}, Path: "", SourceFullPath: ""},
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

func TestFilterStartedWith(t *testing.T) {
	files := []LocalFileInfo{
		{FileInfo: mockFileInfo{name: "prefix1_sample.txt"}, Path: "", SourceFullPath: ""},
		{FileInfo: mockFileInfo{name: "prefix2_sample.txt"}, Path: "", SourceFullPath: ""},
		{FileInfo: mockFileInfo{name: "non_bank_file.txt"}, Path: "", SourceFullPath: ""},
	}
	prefixes := []string{"prefix1", "prefix2"}

	filteredFiles := FilterStartedWith(files, prefixes)

	if len(filteredFiles) != 2 {
		t.Errorf("Expected 2 files to match prefixes, got %d", len(filteredFiles))
	}

	expectedNames := map[string]bool{"prefix1_sample.txt": true, "prefix2_sample.txt": true}
	for _, file := range filteredFiles {
		if _, ok := expectedNames[file.Name()]; !ok {
			t.Errorf("Unexpected file: %s", file.Name())
		}
	}
	for _, file := range filteredFiles {
		t.Logf("File: %s", file.Name())
	}
}

func TestFilterAfterDate(t *testing.T) {
	afterDate, _ := time.Parse("02012006", "12012024")

	files := []LocalFileInfo{
		{FileInfo: mockFileInfo{name: "POS_RevAuthFile_16012024.000002"}},
		{FileInfo: mockFileInfo{name: "POS_RevAuthFile_11012024.000003"}},
		{FileInfo: mockFileInfo{name: "POS_RevAuthFile_02012024.000004"}},
		{FileInfo: mockFileInfo{name: "CL.000005.240113"}},
		{FileInfo: mockFileInfo{name: "random_file.txt"}},
	}

	filteredFiles := FilterAfterDate(files, afterDate)

	expectedFileNames := map[string]bool{
		"POS_RevAuthFile_16012024.000002": true,
		"CL.000005.240110":                true,
	}

	if len(filteredFiles) != len(expectedFileNames) {
		t.Errorf("Expected %d files after filtering, got %d", len(expectedFileNames), len(filteredFiles))
	}

	for _, file := range filteredFiles {
		if _, ok := expectedFileNames[file.Name()]; !ok {
			t.Errorf("Unexpected file after filtering: %s", file.Name())
		}
	}

	for _, file := range filteredFiles {
		t.Logf("Filtered File: %s", file.Name())
	}
}
