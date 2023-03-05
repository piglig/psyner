package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// connect to server
	conn, err := net.Dial("tcp", "host.docker.internal:7777")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localDir := "./data"
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	//go func() {
	// TODO periodically check if local directory is consistent with server
	for range ticker.C {
		checkSum := make(map[string]string)
		err = filepath.Walk(localDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			checksum, err := generateChecksum(path)
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(localDir, path)
			if err != nil {
				return err
			}
			checkSum[relPath] = checksum
			fmt.Printf("time:%v %s: %s\n", time.Now(), path, checksum)
			return nil
		})

		if err != nil {
			return
		}

		// TODO compare with server checksum, get not exist file from server
	}
	//}()

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
