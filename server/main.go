package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"psyner/common"
	"psyner/server/taskrun/action"
	_ "psyner/server/taskrun/action/event"
	"sync"
)

const (
	listenAddr = "localhost:8888"
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
	connPool := make(map[string]net.Conn)
	connPoolLock := sync.Mutex{}

	closeCh := make(chan string)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			// add connection to pool
			connPoolLock.Lock()
			connPool[conn.RemoteAddr().String()] = conn
			connPoolLock.Unlock()
			go connectionHandler(conn, closeCh)
		}
	}()

	for {
		select {
		case event := <-fw.Events:
			if event.Has(fsnotify.Write) {
				log.Printf("File %s modified\n", event.Name)

				// transfer updated file to remote computers
				fileName := filepath.Base(event.Name)
				err := transferFile(fileName, dir, &connPool, &connPoolLock)
				if err != nil {
					log.Println(err)
				}
				log.Println(fileName)
			}
		case closeConn := <-closeCh:
			connPoolLock.Lock()
			delete(connPool, closeConn)
			connPoolLock.Unlock()
			log.Printf("Close connection %s......\n", closeConn)
		}
	}
}

func connectionHandler(conn net.Conn, closeCh chan string) {
	defer conn.Close()
	log.Printf("Accept connection from %s......\n", conn.RemoteAddr())
	payload := common.FileSyncPayload{}
	decoder := gob.NewDecoder(conn)
	for {
		err := decoder.Decode(&payload)
		if err != nil {
			log.Println("connectionHandler", err)
			break
		}

		log.Println("Received data:", payload.ActionType, string(payload.ActionPayload))
		err = action.FileSyncAction(context.Background(), payload.ActionType, string(payload.ActionPayload))
		if err != nil {
			log.Printf("connectionHandler err:%s\n", err.Error())
		}
	}

	closeCh <- conn.RemoteAddr().String()

	//fileName = filepath.Join(dir, fileName)
	//file, err := os.Create(fileName)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//defer file.Close()
	//
	//_, err = io.Copy(file, conn)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

}

func transferFile(fileName, folder string, connPool *map[string]net.Conn, connPoolLock *sync.Mutex) error {
	filePath := filepath.Join(folder, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// make a copy of the connection pool to avoid holding the lock for too long
	connPoolLock.Lock()
	poolCopy := make(map[string]net.Conn, len(*connPool))
	for k, v := range *connPool {
		poolCopy[k] = v
	}
	connPoolLock.Unlock()

	// send file to each remote computer
	for _, conn := range poolCopy {
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
