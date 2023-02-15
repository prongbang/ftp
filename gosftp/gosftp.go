package gosftp

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Run() {
	// Connect to the remote server
	config := &ssh.ClientConfig{
		User: "username",
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", "sftp.example.com:22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer conn.Close()

	// Create an SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal("Failed to create client: ", err)
	}
	defer client.Close()

	// Open the remote file
	remoteFile, err := client.Open("/path/to/remote/file.txt")
	if err != nil {
		log.Fatal("Failed to open remote file: ", err)
	}
	defer remoteFile.Close()

	// Create the local file
	localFile, err := os.Create("/path/to/local/file.txt")
	if err != nil {
		log.Fatal("Failed to create local file: ", err)
	}
	defer localFile.Close()

	// Download the file
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		log.Fatal("Failed to download file: ", err)
	}

	fmt.Println("File downloaded successfully.")
}
