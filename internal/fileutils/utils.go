package fileutils

import (
	"os"
	"path/filepath"
	"strings"
	"tt-copier/internal/sftp"
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

func LoadAllSourceFiles(client *sftp.Client, paths []string) ([]LocalFileInfo, error) {
	var files []LocalFileInfo
	for _, path := range paths {
		dirFiles, err := client.ListFiles(path)

		if err != nil {
			return nil, err
		}

		for _, file := range dirFiles {
			if !file.IsDir() {
				fullPath := filepath.Join(path, file.Name())

				files = append(files, LocalFileInfo{
					FileInfo:       file,
					Path:           path,
					SourceFullPath: fullPath,
				})
			}
		}
	}
	return files, nil
}
