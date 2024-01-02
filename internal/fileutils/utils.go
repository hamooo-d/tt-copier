package fileutils

import (
	"os"
	"path/filepath"
	"strings"
)

type FileInfoExtended struct {
	os.FileInfo
	DestinationPath     string
	DestinationFullPath string
	SourceFullPath      string
}

type LocalFileInfo struct {
	os.FileInfo
	Path           string
	SourceFullPath string
}

func AddBankDestination(source []LocalFileInfo, basePath string) ([]FileInfoExtended, error) {
	banksNames := map[string]string{
		"000003": "SB",
		"000002": "ATIB",
		"000004": "NAB",
		"000005": "MED",
		"000001": "TT",
		"000006": "NCB",
	}

	var newFileList []FileInfoExtended

	for _, file := range source {
		for id, name := range banksNames {
			if strings.Contains(file.Name(), id) { // Check if the filename contains the bank ID
				destination := filepath.Join(basePath, name, "Prod", "from_tadawul")
				newFileList = append(newFileList, FileInfoExtended{
					FileInfo:            file,
					DestinationPath:     destination,
					DestinationFullPath: filepath.Join(destination, file.Name()),
				})
				break // Assuming only one bank ID matches per file
			}
		}
	}

	return newFileList, nil
}

func AddGetDestination(fileList []LocalFileInfo, destination string) ([]FileInfoExtended, error) {
	var updatedFileList []FileInfoExtended

	for _, file := range fileList {
		updatedFileList = append(updatedFileList, FileInfoExtended{
			FileInfo:            file,
			DestinationPath:     destination,
			DestinationFullPath: filepath.Join(destination, file.Name()),
		})
	}

	return updatedFileList, nil
}

func FilterStartedWith(files []LocalFileInfo, prefixes []string) []LocalFileInfo {
	var filteredFiles []LocalFileInfo

	for _, file := range files {
		for _, prefix := range prefixes {
			if strings.HasPrefix(file.Name(), prefix) {
				filteredFiles = append(filteredFiles, file)
			}
		}
	}
	return filteredFiles
}

func AddTTDestination(source []LocalFileInfo) ([]FileInfoExtended, error) {
	ttDestination := "/home/sftp/files/TTP/TT/Prod/from_tadawul"

	return AddGetDestination(source, ttDestination)
}

func LoadAllLocalFiles(paths []string) ([]LocalFileInfo, error) {
	var files []LocalFileInfo
	for _, path := range paths {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				fileInfo, err := entry.Info()
				if err != nil {
					return nil, err
				}
				files = append(files, LocalFileInfo{
					FileInfo:       fileInfo,
					Path:           path,
					SourceFullPath: filepath.Join(path, fileInfo.Name()),
				})
			}
		}
	}
	return files, nil
}
