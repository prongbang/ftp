package goftp

import (
	"bytes"
	"log"

	"github.com/jlaffaye/ftp"
)

func Run() {
	// Connect to the FTP server
	conn, err := ftp.Dial("localhost:2121")
	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}
	defer conn.Quit()

	// Login with the provided credentials
	err = conn.Login("test", "test")
	if err != nil {
		log.Fatal("Failed to login: ", err)
	}

	data := bytes.NewBufferString("Hello World")
	err = conn.Stor("file.txt", data)
	if err != nil {
		panic(err)
	}

	// Open the remote file
	// remoteFile, err := conn.Retr("/path/to/remote/file.txt")
	// if err != nil {
	// 	log.Fatal("Failed to open remote file: ", err)
	// }
	// defer remoteFile.Close()

	// Create the local file
	// localFile, err := os.Create("/path/to/local/file.txt")
	// if err != nil {
	// 	log.Fatal("Failed to create local file: ", err)
	// }
	// defer localFile.Close()

	// Download the file
	// buff := []byte{}
	// _, err = remoteFile.Read(buff)
	// if err != nil {
	// 	log.Fatal("Failed to download file: ", err)
	// }

	// fmt.Println("File downloaded successfully.")
}
