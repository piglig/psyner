package main

import (
	"encoding/gob"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
)

const (
	listenAddr = "localhost:7777"
	dir        = "./data"
)

func main() {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer fw.Close()
	if err = fw.Add(dir); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("listening on %s......\n", listenAddr)
	var m sync.Map
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			m.Store(conn.RemoteAddr().String(), conn)
			go connectionHandler(conn)
		}
	}()

	for {
		select {
		case event := <-fw.Events:
			if event.Has(fsnotify.Write) {
				log.Printf("File %s modified\n", event.Name)

				// transfer updated file to remote computers
				fileName := filepath.Base(event.Name)
				err := transferFile(fileName, dir, &connPool)
				if err != nil {
					fmt.Println(err)
				}
				log.Println(fileName)
			}
		}
	}
}

func connectionHandler(conn net.Conn) {
	defer conn.Close()
	log.Printf("Accept connection from %s......\n", conn.RemoteAddr())

	var fileName string
	err := gob.NewDecoder(conn).Decode(&fileName)
	if err != nil {
		log.Println(err)
		return
	}

	fileName = filepath.Join(dir, fileName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, conn)
	if err != nil {
		log.Println(err)
		return
	}
}

func transferFile(fileName, folder string, connPool *sync.Map) error {
	filePath := filepath.Join(folder, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()


	// send file to each remote computer
	for _, conn := range connPool. {
		go func(conn net.Conn) {
			defer conn.Close()

			fmt.Printf("Sending file %s to %s\n", fileName, conn.RemoteAddr())

			// send file name to remote computer
			err := gob.NewEncoder(conn).Encode(fileName)
			if err != nil {
				fmt.Println(err)
				return
			}

			// send file data to remote computer
			_, err = io.Copy(conn, file)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Sent file %s to %s\n", fileName, conn.RemoteAddr())
		}(conn)
	}

	return nil
}
