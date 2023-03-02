package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// connect to server
	conn, err := net.Dial("tcp", "localhost:7777")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	go func() {
		// TODO periodically check if local directory is consistent with server
		for range ticker.C {

		}
	}()

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
