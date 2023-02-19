package gosftp

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	Host, User, Password string
	Port                 int
	*sftp.Client
	SSHTrustedKey string
}

type SftpClient interface {
	PutFile(localFile, remoteFile string) (err error)
	PutString(text, remoteFile string) (err error)
}

type sftpClient struct {
	*sftp.Client
}

// New Create a new SFTP connection by given parameters
func New(c *Config) (SftpClient, error) {
	client, err := c.connect()
	return &sftpClient{Client: client}, err
}

// PutFile Upload file to sftp server
func (sc *sftpClient) PutFile(localFile, remoteFile string) (err error) {
	srcFile, err := os.Open(localFile)
	if err != nil {
		fmt.Println("Open file", err)
		return err
	}
	defer func(srcFile *os.File) { _ = srcFile.Close() }(srcFile)

	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		_ = sc.Mkdir(path)
	}

	dstFile, err := sc.Create(remoteFile)
	if err != nil {
		fmt.Println("Create remote file", err)
		return err
	}
	defer func(dstFile *sftp.File) { _ = dstFile.Close() }(dstFile)

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// PutString Upload message to sftp server
func (sc *sftpClient) PutString(text, remoteFile string) (err error) {
	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		_ = sc.Mkdir(path)
	}

	// Open a file on the remote SFTP server for writing
	dstFile, err := sc.Create(remoteFile)
	if err != nil {
		fmt.Println("Open a file on the remote:", err)
	}
	defer func(dstFile *sftp.File) { _ = dstFile.Close() }(dstFile)

	// Write the string text to the remote file
	_, err = io.Copy(dstFile, strings.NewReader(text))
	if err != nil {
		fmt.Println("Copy a text to remote:", err)
	}
	return nil
}

func (sc *Config) connect() (client *sftp.Client, err error) {
	config := &ssh.ClientConfig{
		User:    sc.User,
		Auth:    []ssh.AuthMethod{ssh.Password(sc.Password)},
		Timeout: 30 * time.Second,
	}
	if sc.SSHTrustedKey != "" {
		config.HostKeyCallback = trustedHostKeyCallback(sc.SSHTrustedKey)
	} else {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	// Connect to ssh
	addr := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Print("Connect fail:", err)
		return nil, err
	}

	// Create sftp client
	client, err = sftp.NewClient(conn)
	if err != nil {
		fmt.Print("New client fail:", err)
		return nil, err
	}

	return client, nil
}

// SSH Key-strings
func trustedHostKeyCallback(trustedKey string) ssh.HostKeyCallback {

	if trustedKey == "" {
		return func(_ string, _ net.Addr, k ssh.PublicKey) error {
			log.Printf("WARNING: SSH-key verification is *NOT* in effect: to fix, add this trustedKey: %q", keyString(k))
			return nil
		}
	}

	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
		ks := keyString(k)
		if trustedKey != ks {
			return fmt.Errorf("SSH-key verification: expected %q but got %q", trustedKey, ks)
		}

		return nil
	}
}

func keyString(k ssh.PublicKey) string {
	return k.Type() + " " + base64.StdEncoding.EncodeToString(k.Marshal())
}
