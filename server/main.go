package main

import (
	"encoding/gob"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

const (
	listenAddr = "localhost:8000"
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
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go connectionHandler(conn)
		}
	}()

	for {
		select {
		case event := <-fw.Events:
			if event.Has(fsnotify.Write) {

				fileName := filepath.Base(event.Name)
				log.Println(fileName)
			}
			//default:
			//	log.Println("default")
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
