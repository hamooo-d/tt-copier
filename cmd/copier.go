package main

import (
	"fmt"
	"os"

	"tt-copier/internal/db"
	"tt-copier/internal/fileutils"
	"tt-copier/internal/logger"
	"tt-copier/internal/sftp"
)

func uploadToSFTP(client *sftp.Client) bool {
	dbInstance, err := db.NewDBInstance("./db.sqlite")

	if err != nil {
		logger.Error("Error creating DB instance: %v", err)

		return false
	}

	logger.Info("DB instance created successfully.", "INIT", "SUCCESS")

	sourceList := []string{
		"/online/mxpprod/selectsystem_files/cardholder/out",
		"/online/mxpprod/selectsystem_files/transaction/out",
		"/online/mxpprod/selectsystem_files/merchant/out",
		"/online/mxpprod/selectsystem_files/evoucher/out",
		"/online/mxpprod/selectsystem_files/merchant/out",
	}

	bankFilesPrefixes := []string{"CL.", "KYCFile_", "reload_",
		"Rev_Reload_", "redemp_", "POS_RevAuthFile_", "APPLICATION.", "KYC_ATM_",
		"SETT_TOPUP.", "CORP_TOPUP.", "EV_MERC"}

	TTFilesPrefixes := []string{"PersoFile", "CL_TT", "FSO.", "PAYOUT."}

	logger.Info("Loading all files from source list.", "LOAD", "START")

	allFiles, err := fileutils.LoadAllLocalFiles(sourceList)

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

	bankFiles := fileutils.FilterStartedWith(filteredFiles, bankFilesPrefixes)
	TTFiles := fileutils.FilterStartedWith(filteredFiles, TTFilesPrefixes)

	logger.Info(fmt.Sprintf("Filtered, remaining %d bank files and %d TT files.", len(bankFiles), len(TTFiles)), "FILTER", "SUCCESS")

	bankFilesWithDestination, err := fileutils.AddBankDestination(bankFiles, "/home/sftp/files/")

	if err != nil {
		logger.Error("Error adding bank destination: %v", err)
		return false
	}

	TTFilesWithDestination, err := fileutils.AddTTDestination(TTFiles) // Implement AddTTDestination

	if err != nil {
		logger.Error("Error adding bank destination: %v", err)
		return false
	}

	for _, file := range bankFilesWithDestination {
		sourcePath := file.SourceFullPath
		destinationPath := file.DestinationFullPath

		logger.Info(fmt.Sprintf("Uploading bank file %s", file.Name()), "UPLOAD", "START")

		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			logger.Warn(fmt.Sprintf("File %s does not exist, skipping.", file.Name()), "UPLOAD", "SKIPPED")
			continue
		}

		err := client.PutProcedure(sourcePath, destinationPath)

		if err != nil {
			logger.Warn(fmt.Sprintf("Error uploading bank file %s: %v", file.Name(), err), "UPLOAD", "FAILED")
		} else {
			logger.Info(fmt.Sprintf("Successfully uploaded bank file %s", file.Name()), "UPLOAD", "SUCCESS")
			dbInstance.LogEntry(sourcePath, destinationPath, file.Name(), "UPLOAD")
		}
	}

	for _, file := range TTFilesWithDestination {
		sourcePath := file.SourceFullPath
		destinationPath := file.DestinationFullPath

		err := client.PutProcedure(sourcePath, destinationPath)
		if err != nil {
			logger.Warn(fmt.Sprintf("Error uploading TT file %s: %v", file.Name(), err), "UPLOAD", "FAILED")
		} else {
			logger.Info(fmt.Sprintf("Successfully uploaded TT file %s", file.Name()), "UPLOAD", "SUCCESS")
			dbInstance.LogEntry(sourcePath, destinationPath, file.Name(), "UPLOAD")
		}
	}

	return true
}

func main() {
	logger.Init()
	client, err := sftp.NewClient("192.168.100.6", 22, "sofian", "sofian")

	if err != nil {
		logger.Info("Error creating SFTP client, exiting.", "UPLOAD", "FAILED")
		return
	}

	defer client.Close()

	success := uploadToSFTP(client)

	if !success {
		logger.Info("Upload failed, exiting.", "UPLOAD", "FAILED")
		os.Exit(1)
	} else {
		logger.Info("Upload finished successfully.", "UPLOAD", "SUCCESS")
	}
}
