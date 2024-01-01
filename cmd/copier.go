package main

import (
	"log"
	"tt-copier/internal/sftp"
)

func main() {
	client, err := sftp.NewClient("192.168.100.6", 22, "atib", "123")

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	files, err := client.ListFiles("./")

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		log.Println(file.Name())
	}
}
