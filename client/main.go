package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"psyner/client/cmd"
	"strings"
	"time"
)

func main() {
	// connect to server
	// host.docker.internal
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localDir := "./data"
	client, err := cmd.NewClient(cmd.ClientConfig{
		LocalDir:       localDir,
		TickerInterval: 10 * time.Second,
	})
	if err != nil {
		log.Fatal("NewClient", err.Error())
	}

	client.Start()

	fmt.Println("read file...")

	for {
		// read available files from server
		var fileName string
		err = gob.NewDecoder(conn).Decode(&fileName)
		if err != nil {
			log.Fatal(err)
		}
		fileName = strings.TrimSpace(fileName)

		// send selected file name to server
		err = gob.NewEncoder(conn).Encode(fileName)
		if err != nil {
			log.Fatal(err)
		}

		// receive file data from server
		fileData := bytes.Buffer{}
		_, err = io.Copy(&fileData, conn)
		if err != nil {
			log.Fatal(err)
		}

		// save file to local computer
		err = os.WriteFile(fileName, fileData.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Saved file %s\n", fileName)
	}
}

func generateChecksum(filePath string) (string, error) {
	// Generate the checksum for a file given its path
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
