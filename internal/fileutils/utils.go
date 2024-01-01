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
}

func AddBankDestination(source []os.FileInfo, basePath string) ([]FileInfoExtended, error) {
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

func AddGetDestination(fileList []os.FileInfo, destination string) ([]FileInfoExtended, error) {
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
