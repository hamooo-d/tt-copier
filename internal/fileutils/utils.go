package fileutils

import (
	"os"
	"path/filepath"
	"strings"
	"time"
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

func AddBankDestination(source []LocalFileInfo, basePath string, banksNames map[string]string, env string) ([]FileInfoExtended, error) {
	var newFileList []FileInfoExtended

	for _, file := range source {
		for id, name := range banksNames {
			if strings.Contains(file.Name(), id) {
				destination := filepath.Join(basePath, name, env, "from_tadawul")
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
		for _, prefix := range prefixes {
			if strings.HasPrefix(file.Name(), prefix) {
				filteredFiles = append(filteredFiles, file)
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

func FilterAfterDate(files []LocalFileInfo, afterDate time.Time) []LocalFileInfo {
	var filteredFiles []LocalFileInfo

	const layout = "02012006"
	const dateLength = 8

	for _, file := range files {
		name := file.Name()

		var dateStr string

		for i := len(name) - dateLength; i >= 0; i-- {
			substr := name[i : i+dateLength]

			if _, err := time.Parse(layout, substr); err == nil {
				dateStr = substr
				break
			}
		}

		if dateStr == "" {
			continue
		}

		fileDate, err := time.Parse(layout, dateStr)

		if err != nil {
			continue
		}

		if fileDate.After(afterDate) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}

func LoadAllSourceFiles(paths []string) ([]LocalFileInfo, error) {
	var files []LocalFileInfo
	for _, path := range paths {
		dirFiles, err := os.ReadDir(path) // Use ioutil.ReadDir for older Go versions
		if err != nil {
			return nil, err
		}

		for _, dirEntry := range dirFiles {
			if dirEntry.IsDir() {
				continue
			}

			fileInfo, err := dirEntry.Info()
			if err != nil {
				return nil, err
			}

			fullPath := filepath.Join(path, fileInfo.Name())

			files = append(files, LocalFileInfo{
				FileInfo:       fileInfo,
				Path:           path,
				SourceFullPath: fullPath,
			})
		}
	}
	return files, nil
}
