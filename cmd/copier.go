package main

import (
	"fmt"
	"os"

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

	allFiles, err := fileutils.LoadAllSourceFiles(client, sourceList)

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

	logger.Info(fmt.Sprintf("Filtered, remaining %d files not uploaded.", len(filteredFiles)), "FILTER", "SUCCESS")

	logger.Info("Filtering bank and TT files.", "FILTER", "START")

	bankFiles := fileutils.FilterStartedWith(filteredFiles, bankPrefixes)
	TTFiles := fileutils.FilterStartedWith(filteredFiles, TTPrefixes)

	logger.Info(fmt.Sprintf("Filtered, remaining %d bank files and %d TT files.", len(bankFiles), len(TTFiles)), "FILTER", "SUCCESS")

	bankFilesWithDestination, err := fileutils.AddBankDestination(bankFiles, cfg.Dests.BankDest, cfg.BanksNames, cfg.Env)

	if err != nil {
		logger.Error("Error adding bank destination: %v", err)

		return false
	}

	TTFilesWithDestination, err := fileutils.AddTTDestination(TTFiles)

	if err != nil {
		logger.Error("Error adding bank destination: %v", err)
		return false
	}

	bankUploadCount := 0

	for _, file := range bankFilesWithDestination {
		sourcePath := file.SourceFullPath
		destinationPath := file.DestinationFullPath

		logger.Info(fmt.Sprintf("Uploading bank file %s", file.Name()), "UPLOAD", "START")

		err := client.PutProcedure(sourcePath, destinationPath)

		if err != nil {
			logger.Warn(fmt.Sprintf("Error uploading bank file %s: %v", file.Name(), err), "UPLOAD", "FAILED")
		} else {
			err = dbInstance.LogEntry(sourcePath, destinationPath, file.Name())

			if err != nil {
				logger.Error("Error logging entry: %v", err)
			} else {
				logger.Info("File saved in db", "UPLOAD", "SUCCESS")

			}
			bankUploadCount++
			logger.Info(fmt.Sprintf("Successfully uploaded bank file %s", file.Name()), "UPLOAD", "SUCCESS")
		}
	}

	// log total bank files and uploaded bank files
	logger.Info(fmt.Sprintf("Total bank files: %d", len(bankFiles)), "UPLOAD", "SUCCESS")
	logger.Info(fmt.Sprintf("Uploaded %d bank files, Total", bankUploadCount), "UPLOAD", "SUCCESS")

	for _, file := range TTFilesWithDestination {
		sourcePath := file.SourceFullPath
		destinationPath := file.DestinationFullPath

		err := client.PutProcedure(sourcePath, destinationPath)
		if err != nil {
			logger.Warn(fmt.Sprintf("Error uploading TT file %s: %v", file.Name(), err), "UPLOAD", "FAILED")
		} else {
			logger.Info(fmt.Sprintf("Successfully uploaded TT file %s", file.Name()), "UPLOAD", "SUCCESS")
		}
	}

	return true
}

func main() {
	logger.Init()
	cfg, err := config.LoadConfig(".")

	if err != nil {
		logger.Error("Error loading config: %v", err)
		os.Exit(1)
	}

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
		logger.Info("Upload finished successfully.", "UPLOAD", "SUCCESS")
	}

}
