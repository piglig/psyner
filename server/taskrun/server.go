package taskrun

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"io/fs"
	"log"
	"net"
	"os"
	"path/filepath"
	"psyner/common"
	"sync"
	"time"
)

type Server struct {
	checkSumMux     sync.RWMutex
	connPoolMux     sync.RWMutex
	relPathCheckSum map[string]string
	config          ServerConfig
	connPool        map[string]net.Conn
}

type ServerConfig struct {
	ListenAddr string
	LocalDir   string
}

type Context struct {
	*Server
	context.Context
}

func NewServer(config ServerConfig) (*Server, error) {
	if config.LocalDir == "" {
		return nil, fmt.Errorf("local dir %s not invalid", config.LocalDir)
	}

	_, err := os.Stat(config.LocalDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("local dir %s not exist", config.LocalDir)
		} else {
			return nil, fmt.Errorf("local dir stat invalid %v", err)
		}
	}

	return &Server{
		connPool:        make(map[string]net.Conn),
		relPathCheckSum: make(map[string]string),
		config:          config,
	}, err
}

func (s *Server) checkLocalDirChecksum(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := filepath.Walk(s.config.LocalDir, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.Mode().IsRegular() {
					return nil
				}

				calSum, err := common.GenerateChecksum(path)
				if err != nil {
					return err
				}

				relPath, err := filepath.Rel(s.config.LocalDir, path)
				if err != nil {
					return err
				}

				checkSum, ok := s.relPathCheckSum[relPath]
				if ok && checkSum == calSum {
					return nil
				}

				s.checkSumMux.Lock()
				s.relPathCheckSum[relPath] = calSum
				s.checkSumMux.Unlock()
				log.Printf("%s: %s\n", path, calSum)
				return nil
			})

			if err != nil {
				log.Println("checkLocalDirChecksum", "err", err.Error())
				return
			}
		}
	}
}

func (s *Server) Run() {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer fw.Close()
	if err = fw.Add(s.config.LocalDir); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", s.config.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	go s.checkLocalDirChecksum(5 * time.Second)
	go s.connectionHandler(listener)

	for {
		select {
		case event := <-fw.Events:
			if event.Has(fsnotify.Write) {
				log.Printf("File %s modified\n", event.Name)
				// transfer updated file to remote computers
				fileName := filepath.Base(event.Name)
				//err := transferFile(fileName, dir, &connPool, &connPoolLock)
				//if err != nil {
				//	log.Println(err)
				//}
				log.Println(fileName)
			}
		case err, ok := <-fw.Errors:
			if !ok {
				return
			}
			log.Println("fw error", err.Error())
		}
	}
}

func (s *Server) connectionHandler(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		s.connPoolMux.Lock()
		s.connPool[conn.RemoteAddr().String()] = conn
		s.connPoolMux.Unlock()

		go func(conn net.Conn) {
			defer func() {
				conn.Close()
				s.connPoolMux.Lock()
				delete(s.connPool, conn.RemoteAddr().String())
				s.connPoolMux.Unlock()
				log.Printf("Close connection %s......\n", conn.RemoteAddr().String())
			}()
			log.Printf("Accept connection from %s......\n", conn.RemoteAddr())
			payload := common.FileSyncPayload{}
			decoder := gob.NewDecoder(conn)
			for {
				err := decoder.Decode(&payload)
				if err != nil {
					log.Println("connectionHandler", err)
					break
				}

				//log.Println("Received data:", payload.ActionType, string(payload.ActionPayload))
				//ctx := Context{
				//	Server:  s,
				//	Context: context.Background(),
				//}
				//err = FileSyncAction(ctx, payload.ActionType, conn, string(payload.ActionPayload))
				//if err != nil {
				//	log.Printf("connectionHandler err:%s\n", err.Error())
				//}
			}

		}(conn)
	}
}

func (s *Server) CheckFileExist(path string) bool {
	s.checkSumMux.RLock()
	defer s.checkSumMux.RUnlock()
	_, ex := s.relPathCheckSum[path]
	return ex
}

func (s *Server) getFilePath(path string) string {
	return filepath.Join(s.config.LocalDir, path)
}

func FileSyncAction(ctx Context, action common.FileSyncOp, conn net.Conn, command string) error {
	return Exec(ctx, action, conn, command)
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

func (s *Server) broadcastAction() {
	s.connPoolMux.Lock()
	poolCopy := make(map[string]net.Conn, len(s.connPool))
	for k, v := range s.connPool {
		poolCopy[k] = v
	}
	s.connPoolMux.Unlock()

	//for _, conn := range poolCopy {
	//
	//}
}
