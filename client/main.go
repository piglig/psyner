package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		// read available files from server
		var files []string
		err = gob.NewDecoder(conn).Decode(&files)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Available files:")
		for _, file := range files {
			fmt.Println(file)
		}

		// select file to synchronize
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter file to synchronize: ")
		fileName, err := reader.ReadString('\n')
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
