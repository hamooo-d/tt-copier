package sftp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	sftpClient *sftp.Client
}

func NewClient(host string, port int, username, password string) (*Client, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %s", err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, fmt.Errorf("failed to start SFTP client: %s", err)
	}

	return &Client{sftpClient: sftpClient}, nil
}

func (c *Client) Close() {
	c.sftpClient.Close()
}

func (c *Client) ListFiles(dir string) ([]os.FileInfo, error) {
	return c.sftpClient.ReadDir(dir)
}

func (c *Client) CopyRename(files []os.FileInfo, source, destination string, renameRules map[string]string) error {
	for _, file := range files {
		newFilename := file.Name()

		for initial, newPrefix := range renameRules {
			if strings.Contains(file.Name(), initial) {
				newFilename = strings.ReplaceAll(file.Name(), initial, newPrefix)
				break
			}
		}
		sourcePath := filepath.Join(source, file.Name())
		destinationPath := filepath.Join(destination, newFilename)

		inFile, err := c.sftpClient.Open(sourcePath)
		if err != nil {
			return err
		}
		defer inFile.Close()

		outFile, err := c.sftpClient.Create(destinationPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, inFile); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) PutProcedure(localPath, remotePath string) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	remoteFile, err := c.sftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, localFile)
	return err
}

func (c *Client) GetProcedure(remotePath, localPath string) error {
	remoteFile, err := c.sftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, remoteFile)
	return err
}
