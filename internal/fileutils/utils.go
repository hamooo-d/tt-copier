package fileutils

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
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
func AddTTDestination(source []LocalFileInfo) ([]FileInfoExtended, error) {
	ttDestination := "/home/sftp/files/TT/Prod/from_tadawul"

	return AddGetDestination(source, ttDestination)
}

func FilterAfterDate(files []LocalFileInfo, afterDate time.Time) []LocalFileInfo {
	var filteredFiles []LocalFileInfo

	var dateRegexes = []*regexp.Regexp{
		regexp.MustCompile(`\d{8}`),
		regexp.MustCompile(`\d{6}`),
	}

	for _, file := range files {
		name := file.Name()
		dateStr := extractDateFromName(name, dateRegexes)

		if dateStr == "" {
			continue
		}

		fileDate, err := parseDate(dateStr)
		if err != nil {
			continue
		}

		if fileDate.After(afterDate) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}

func extractDateFromName(name string, regexes []*regexp.Regexp) string {
	if strings.HasPrefix(name, "PersoFile_") {
		dotIndex := strings.LastIndex(name, ".")
		if dotIndex != -1 && dotIndex >= 6 {
			datePart := name[dotIndex-6 : dotIndex]
			if _, err := time.Parse("060102", datePart); err == nil {
				return datePart
			}
		}
	}

	for _, regex := range regexes {
		matches := regex.FindAllString(name, -1)
		if matches != nil {
			return matches[len(matches)-1]
		}
	}
	return ""
}

func parseDate(dateStr string) (time.Time, error) {
	var layouts = []string{
		"02012006",
		"060102",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("unable to parse date")
}

func LoadAllSourceFiles(paths []string) ([]LocalFileInfo, error) {
	var files []LocalFileInfo
	for _, path := range paths {
		dirFiles, err := os.ReadDir(path)
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
