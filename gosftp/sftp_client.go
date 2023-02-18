package gosftp

import (
	"encoding/base64"
	"errors"
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
	SSHTurstedKey string
}

// Create a new SFTP connection by given parameters
func New(c *Config) (client *Config, err error) {
	switch {
	case `` == strings.TrimSpace(c.Host),
		`` == strings.TrimSpace(c.User),
		`` == strings.TrimSpace(c.Password),
		0 >= c.Port || c.Port > 65535:
		return nil, errors.New("Invalid parameters")
	}

	if err = c.connect(); nil != err {
		return nil, err
	}
	return c, nil
}

// Upload file to sftp server
func (sc *Config) PutFile(localFile, remoteFile string) (err error) {
	srcFile, err := os.Open(localFile)
	if err != nil {
		fmt.Println("Open file", err)
		return err
	}
	defer srcFile.Close()

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
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// Upload message to sftp server
func (sc *Config) PutMessage(message, remoteFile string) (err error) {
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
	defer dstFile.Close()

	// Write the string message to the remote file
	_, err = io.Copy(dstFile, strings.NewReader(message))
	if err != nil {
		fmt.Println("Copy a message to remote:", err)
	}
	return nil
}

func (sc *Config) connect() (err error) {
	config := &ssh.ClientConfig{
		User:    sc.User,
		Auth:    []ssh.AuthMethod{ssh.Password(sc.Password)},
		Timeout: 30 * time.Second,
	}
	if sc.SSHTurstedKey != "" {
		config.HostKeyCallback = trustedHostKeyCallback(sc.SSHTurstedKey)
	} else {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	// Connet to ssh
	addr := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Print("Connect fail:", err)
		return err
	}

	// Create sftp client
	client, err := sftp.NewClient(conn, sftp.UseFstat(true), sftp.MaxPacket(20480))
	if err != nil {
		fmt.Print("New client fail:", err)
		return err
	}
	sc.Client = client

	return nil
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
