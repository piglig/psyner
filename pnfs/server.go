package pnfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"sync"
	"utils"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type NFSServer interface {
	PING() string
	PostFileTo(string, string)
	GetFileList(w http.ResponseWriter, r *http.Request)
	ReceiveFileFrom()

	// use for http server
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// pnfs file struct
type serverFile struct {
	fileName string
	fileInfo os.FileInfo
	md5      string
}

type PServer struct {
	files      map[string]map[string]serverFile // the other server node and files md5 string
	nodes      []string                         // the server node host
	filePath   string                           // the server node file path
	localFiles []serverFile                     // the current server node files

	addr string // ther server node host and port

	mu sync.Mutex // protects currently request
}

// New initial pnfs server
func New(addr, filePath string, nodes []string) *PServer {
	fmt.Printf("addr [%s], local file path[%s], server nodes%v\n", addr, filePath, nodes)
	return &PServer{
		addr:       addr,
		files:      make(map[string]map[string]serverFile),
		nodes:      nodes,
		filePath:   filePath,
		localFiles: getPathFiles(filePath),
	}
}

func getPathFiles(filePath string) []serverFile {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		log.Fatalf("read path[%s] files error: %v", filePath, err)
		return nil
	}

	serverFiles := []serverFile{}
	for _, file := range files {
		serverFile := serverFile{}
		// serverFile.fileInfo = file
		serverFile.fileName = file.Name()
		serverFile.md5 = utils.MD5(file.Name())
		serverFiles = append(serverFiles, serverFile)
	}
	return serverFiles
}

type FileListReq struct {
	Host     string   `json:"host"`
	FileList []string `json:"file_list"`
}

func PostLocalFiles(host, api, filePath string) {
	res := &FileListReq{}
	localFiles := getPathFiles(filePath)
	localFilesSlice := []string{}
	for _, file := range localFiles {
		localFilesSlice = append(localFilesSlice, file.fileName)
	}

	res.FileList = localFilesSlice
	res.Host = host

	jsonData, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
		return
	}

	resp, err := http.Post(host+api, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("postLocalFiles resp code:%v\n", resp.StatusCode)
		return
	}

	respStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	if string(respStr) != SUCCESS {
		log.Printf("postLocalFiles resp str:%s", string(respStr))
	}
}

const (
	SUCCESS = "success"
	FAIL    = "fail"
)

func (s *PServer) GetFileList(w http.ResponseWriter, r *http.Request) {
	var req FileListReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if len(req.Host) > 0 && len(req.FileList) > 0 {
		serverFiles := map[string]serverFile{}
		for _, file := range req.FileList {
			serverFile := serverFile{}
			serverFile.fileName = file
			serverFile.md5 = utils.MD5(file)
			serverFiles[file] = serverFile
		}
		s.files[req.Host] = serverFiles
	}

	fmt.Printf("%v\n", s.files[req.Host])

	fmt.Fprintf(w, "%s", SUCCESS)
}

func getRemoteFiles(host string) {

}

func (s *PServer) PING() string {
	return "PING"
}

func (s *PServer) PONG(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "PONG")
}

func (s *PServer) PostFileTo(url, fileName string) {
	// file, handler, err := r.FormFile("file")
}

func (s *PServer) SyncWithRemoteNode() {

	for _, localFile := range s.localFiles {
		for host, remoteFile := range s.files {
			if _, ok := remoteFile[localFile.fileName]; !ok {
				s.PostFileTo(host, localFile.fileName)
			}
		}
	}
}

// receive file from remote server node
func (s *PServer) ReceiveFileFrom() {
	resp, err := http.Get("http://10.10.4.54:9998/upload?file=girl.jpg")
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(params)
	fileName := params["filename"]

	out, err := os.Create(s.filePath + "/" + fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()

	// try to save file three times
	for i := 0; i < 3; i++ {
		_, err = io.Copy(out, resp.Body)
		if err == nil {
			break
		}
	}
}

type FileListRes struct {
	Host     string   `json:"host"`
	FileList []string `json:"file_list"`
}

func (s *PServer) PostFileList() {
	var res FileListRes

	if len(s.localFiles) > 0 {
		serverFiles := []string{}
		for _, file := range s.localFiles {
			serverFile := serverFile{}
			serverFile.fileName = file.fileName
			serverFile.md5 = file.md5
			serverFiles = append(serverFiles, serverFile.fileName)
		}
		res.Host = s.addr
		res.FileList = serverFiles
	}

	// fmt.Fprintf(w, "req: %v", req)
}

func (s *PServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/ping":
		s.PONG(w, req)
	case "/getFileList":
		s.GetFileList(w, req)
	}
}

func Run(nfs NFSServer) (err error) {
	s := nfs.(*PServer)
	return http.ListenAndServe(s.addr, s)
}
