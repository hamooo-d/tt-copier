package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"tt-copier/config"
	"tt-copier/internal/db"
	"tt-copier/internal/fileutils"
	"tt-copier/internal/logger"
	"tt-copier/internal/sftp"
)

func uploadToSFTP(client *sftp.Client, cfg *config.Config) bool {
	dbInstance, err := db.NewDBInstance(cfg.Database.DBPath)

	if err != nil {
		logger.Error("Error creating DB instance: %v", err)

		return false
	}

	logger.Info("DB instance created successfully.", "INIT", "SUCCESS")

	sourceList := cfg.SourceList

	bankPrefixes := cfg.FilesPrefixes.BankFilesPrefixes
	TTPrefixes := cfg.FilesPrefixes.TTFilesPrefixes

	logger.Info("Loading all files from source list.", "LOAD", "START")

	allFiles, err := fileutils.LoadAllSourceFiles(sourceList)

	if err != nil {
		logger.Error("Error loading files from source list: %v", err)

		return false
	}

	logger.Info(fmt.Sprintf("Loaded %d files.", len(allFiles)), "LOAD", "SUCCESS")

	logger.Info("Filtering uploaded files.", "FILTER", "START")

	filteredFiles, err := db.FilterUploadedFiles(dbInstance, allFiles)

	if err != nil {
		logger.Error("Error filtering uploaded files: %v", err)

		return false
	}

	if len(filteredFiles) == 0 {
		logger.Info("No files to upload.", "UPLOAD", "SKIPPED")

		return true
	}

	logger.Info("Filtering bank and TT files.", "FILTER", "START")

	bankFiles := fileutils.FilterStartedWith(filteredFiles, bankPrefixes)
	TTFiles := fileutils.FilterStartedWith(filteredFiles, TTPrefixes)

	logger.Info(fmt.Sprintf("Prefixes matched on %d bank files and %d TT files.", len(bankFiles), len(TTFiles)), "FILTER", "SUCCESS")

	afterDate, err := time.Parse("02012006", cfg.AfterDate)

	if err != nil {
		logger.Error("Error parsing after date:", err)

		return false
	}

	bankFiles = fileutils.FilterAfterDate(bankFiles, afterDate)
	TTFiles = fileutils.FilterAfterDate(TTFiles, afterDate)

	logger.Info(fmt.Sprintf("Verified date on %d bank files and %d TT files.", len(bankFiles), len(TTFiles)), "FILTER", "SUCCESS")

	if len(bankFiles) == 0 && len(TTFiles) == 0 {
		logger.Info("No files to upload.", "UPLOAD", "SUCCESS")

		return true
	}

	bankFilesWithDestination, err := fileutils.AddBankDestination(bankFiles, cfg.Dests.BankDest, cfg.BanksNames, cfg.Env)

	logger.Info(fmt.Sprintf("Added destination to %d bank files.", len(bankFilesWithDestination)), "UPLOAD", "SUCCESS")

	if err != nil {
		logger.Error("Error adding bank destination: %v", err)

		return false
	}

	TTFilesWithDestination, err := fileutils.AddTTDestination(TTFiles)

	if err != nil {
		logger.Error("Error adding TT destination: %v", err)
		return false
	}

	logger.Info(fmt.Sprintf("Added destination to %d TT files.", len(TTFilesWithDestination)), "UPLOAD", "SUCCESS")

	bankUploadCount := 0

	var wg sync.WaitGroup
	errChan := make(chan error, len(bankFilesWithDestination))
	semaphore := make(chan struct{}, 10)

	for _, file := range bankFilesWithDestination {
		wg.Add(1)

		go func(file fileutils.FileInfoExtended) {
			defer wg.Done()

			semaphore <- struct{}{}

			sourcePath := file.SourceFullPath
			destinationPath := file.DestinationFullPath

			err := client.PutProcedure(sourcePath, destinationPath)

			if err != nil {
				errChan <- err
			} else {
				bankUploadCount++
				err := dbInstance.LogEntry(sourcePath, destinationPath, file.Name())
				if err != nil {
					logger.Warn(fmt.Sprintf("Error logging file %s: %v", file.Name(), err), "UPLOAD", "FAILED")
				}
			}

			<-semaphore
		}(file)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			logger.Error("Error uploading bank file", err)
		}

	}

	logger.Info(fmt.Sprintf("Total bank files: %d", len(bankFilesWithDestination)), "UPLOAD", "INFORMATIONAL")
	logger.Info(fmt.Sprintf("Uploaded %d bank files, Total", bankUploadCount), "UPLOAD", "INFORMATIONAL")

	logger.Info("Uploading TT files.", "UPLOAD", "START")

	ttUploadCount := 0
	var ttWg sync.WaitGroup
	ttErrChan := make(chan error, len(TTFilesWithDestination))
	ttSemaphore := make(chan struct{}, 10)

	for _, file := range TTFilesWithDestination {
		ttWg.Add(1)

		go func(file fileutils.FileInfoExtended) {
			defer ttWg.Done()

			ttSemaphore <- struct{}{}

			sourcePath := file.SourceFullPath
			destinationPath := file.DestinationFullPath

			err := client.PutProcedure(sourcePath, destinationPath)

			if err != nil {
				ttErrChan <- err
			} else {
				ttUploadCount++
				err := dbInstance.LogEntry(sourcePath, destinationPath, file.Name())
				if err != nil {
					logger.Warn(fmt.Sprintf("Error logging file %s: %v", file.Name(), err), "UPLOAD", "FAILED")
				}
			}

			<-ttSemaphore
		}(file)
	}

	ttWg.Wait()

	close(ttErrChan)

	for err := range ttErrChan {
		if err != nil {
			logger.Error("Error uploading TT file: %v", err)
		}
	}

	logger.Info(fmt.Sprintf("Total TT files: %d", len(TTFilesWithDestination)), "UPLOAD", "INFO")
	logger.Info(fmt.Sprintf("Uploaded %d TT files, Total", ttUploadCount), "UPLOAD", "INFO")

	logger.Info(fmt.Sprintf("Skipped %d files", len(filteredFiles)-(bankUploadCount+ttUploadCount)), "UPLOAD", "INFO")
	logger.Info(fmt.Sprintf("Uploaded %d files, Total", bankUploadCount+ttUploadCount), "UPLOAD", "INFO")

	return true
}

func main() {
	cfg, err := config.LoadConfig(".")

	if err != nil {
		fmt.Println("Error loading config file, exiting. err")
		os.Exit(1)
	}

	logger.Init(cfg)

	client, err := sftp.NewClient(cfg.SFTP.Host, cfg.SFTP.Port, cfg.SFTP.User, cfg.SFTP.Password)

	if err != nil {
		logger.Info("Error creating SFTP client, exiting.", "UPLOAD", "FAILED")
		return
	}

	defer client.Close()

	success := uploadToSFTP(client, cfg)

	if !success {
		logger.Info("Upload failed, exiting.", "UPLOAD", "FAILED")
		os.Exit(1)
	} else {
		logger.Info("Upload finished.", "UPLOAD", "SUCCESS")
	}
}
