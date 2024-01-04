package fileutils

import (
	"log"
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
		"000001": "TT",
		"000002": "ATIB",
		"000003": "SB",
		"000004": "NAB",
		"000005": "MED",
		"000006": "NCB",
	}

	var newFileList []FileInfoExtended

	for _, file := range source {
		for id, name := range banksNames {
			if strings.Contains(file.Name(), id) {
				destination := filepath.Join(basePath, name, "Prod", "from_tadawul")
				newFileList = append(newFileList, FileInfoExtended{
					FileInfo:            file,
					DestinationPath:     destination,
					DestinationFullPath: filepath.Join(destination, file.Name()),
					SourceFullPath:      file.SourceFullPath,
				})
				break
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
			SourceFullPath:      file.SourceFullPath,
		})
	}

	return updatedFileList, nil
}

func FilterStartedWith(files []LocalFileInfo, prefixes []string) []LocalFileInfo {
	var filteredFiles []LocalFileInfo
	for _, file := range files {
		parts := strings.SplitN(file.Name(), "_", 3)

		if len(parts) < 2 {
			continue
		}

		prefixPart := parts[1]

		for _, prefix := range prefixes {
			if strings.HasPrefix(prefixPart, prefix) {
				filteredFiles = append(filteredFiles, file)
				log.Printf("File %s matched prefix %s", file.Name(), prefix)
				break
			}
		}
	}
	return filteredFiles
}

func AddTTDestination(source []LocalFileInfo) ([]FileInfoExtended, error) {
	ttDestination := "/home/sftp/files/TTP/"

	return AddGetDestination(source, ttDestination)
}

func LoadAllLocalFiles(paths []string) ([]LocalFileInfo, error) {
	var files []LocalFileInfo
	for _, path := range paths {
		dirFiles, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		for _, file := range dirFiles {
			if !file.IsDir() {
				fullPath := filepath.Join(path, file.Name())
				fileInfo, err := file.Info()

				if err != nil {
					return nil, err
				}

				files = append(files, LocalFileInfo{
					FileInfo:       fileInfo,
					Path:           path,
					SourceFullPath: fullPath,
				})
			}
		}
	}
	return files, nil
}
