package main

import (
	"fmt"

	"tt-copier/internal/db"
	"tt-copier/internal/fileutils"
	"tt-copier/internal/sftp"
)

func putProcedure(client *sftp.Client) {
	dbInstance, err := db.NewDBInstance("./db.sqlite")

	if err != nil {
		fmt.Printf("Error creating DB instance: %v\n", err)
	}

	sourceList := []string{
		"/online/mxpprod/selectsystem_files/cardholder/in",
		"/online/mxpprod/selectsystem_files/transaction/out",
		"/online/mxpprod/selectsystem_files/merchant/in",
		"/online/mxpprod/selectsystem_files/evoucher/out",
		"/online/mxpprod/selectsystem_files/merchant/out",
	}

	bankFilesPrefixes := []string{"CL.", "KYCFile_", "reload_" /* ... other prefixes */}
	TTFilesPrefixes := []string{"PersoFile", "CL_TT", "FSO.", "PAYOUT."}

	allFiles, err := fileutils.LoadAllLocalFiles(sourceList)
	if err != nil {
		fmt.Printf("Error loading files: %v\n", err)
		return
	}

	filteredFiles, err := db.FilterNotRenamed(dbInstance, allFiles)
	if err != nil {
		fmt.Printf("Error filtering files: %v\n", err)
	}

	filteredFiles, err = db.KeepNotPut(dbInstance, filteredFiles)

	if err != nil {
		fmt.Printf("Error filtering files: %v\n", err)
	}

	// Filter files specific to banks and TT
	bankFiles := fileutils.FilterStartedWith(filteredFiles, bankFilesPrefixes)
	TTFiles := fileutils.FilterStartedWith(filteredFiles, TTFilesPrefixes)

	// Add destination paths
	bankFilesWithDestination, err := fileutils.AddBankDestination(bankFiles, "/path/to/bank/destination")
	if err != nil {
		fmt.Printf("Error adding bank destination: %v\n", err)
		return
	}

	TTFilesWithDestination, err := fileutils.AddTTDestination(TTFiles) // Implement AddTTDestination

	if err != nil {
		fmt.Printf("Error adding TT destination: %v\n", err)
		return
	}

	for _, file := range bankFilesWithDestination {
		sourcePath := file.SourceFullPath
		destinationPath := file.DestinationFullPath

		err := client.PutProcedure(sourcePath, destinationPath)
		if err != nil {
			fmt.Printf("Error uploading file %s: %v\n", file.Name(), err)
		}
	}

	for _, file := range TTFilesWithDestination {
		sourcePath := file.SourceFullPath
		destinationPath := file.DestinationFullPath

		err := client.PutProcedure(sourcePath, destinationPath)
		if err != nil {
			fmt.Printf("Error uploading file %s: %v\n", file.Name(), err)
		}
	}

	fmt.Println("PUT procedure completed successfully.")
}

func main() {
	client, err := sftp.NewClient("sofian", 22, "192.168.100.6", "sofian")
	if err != nil {
		fmt.Printf("Error creating SFTP client: %v\n", err)
		return
	}

	defer client.Close()

	putProcedure(client)
}
