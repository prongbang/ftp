package gosftp

import (
	"fmt"
	"time"
)

func Run() {

	fmt.Println("Sftp Upload started ", time.Now().String())

	config := &Config{
		Host:          "127.0.0.1",
		User:          "atv",
		Password:      "Atv@123!",
		Port:          2022,
		SSHTrustedKey: "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBDyb0S57jLYhQv/qoOoYrK5XBVpEflCVxyboNOaVWFzLu1q6juqBHkdDp6QFvq0ad9x3Kx8qilrl9IhQoktl4fw=",
	}
	ftpClient, err := New(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ftpClient.PutFile("files/test.csv", "localfolder/test.csv")
	err = ftpClient.PutString("1,hello,message", "csv/test.csv")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Sftp Upload finished ", time.Now().String())
	}

}
