package pnfs

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"sync"
	"utils"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type NFSServer interface {
	PING() string
	// server for other server node download
	UploadFileTo(writer http.ResponseWriter, request *http.Request)
	GetLocalFileList(w http.ResponseWriter, r *http.Request)

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

func (s *PServer) PostLocalFileList(host, api, filePath string) {
	// res := &FileListReq{}
	// localFiles := getPathFiles(filePath)
	// localFilesSlice := []string{}
	// for _, file := range localFiles {
	// 	localFilesSlice = append(localFilesSlice, file.fileName)
	// }

	// res.FileList = localFilesSlice
	// res.Host = host

	// jsonData, err := json.Marshal(res)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// resp, err := http.Post(host+api, "application/json", bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	log.Printf("postLocalFiles resp code:%v\n", resp.StatusCode)
	// 	return
	// }

	// respStr, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// }

	// if string(respStr) != SUCCESS {
	// 	log.Printf("postLocalFiles resp str:%s", string(respStr))
	// }
}

const (
	SUCCESS = "success"
	FAIL    = "fail"
)

type LocalFilesRes struct {
	Files []string `json:"files"`
}

func (s *PServer) GetLocalFileList(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := &LocalFilesRes{}
	for _, file := range s.localFiles {
		res.Files = append(res.Files, file.fileName)
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Printf("postLocalFiles marshal to json err:%v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)

	// fmt.Printf("%s request get local[%s] files\n", utils.GetIP(r), s.addr)
}

func (s *PServer) getRemoteFiles(host, api string) {
	addr := "http://" + host + api
	resp, err := http.Get(addr)

	if err != nil {
		log.Printf("%s request get remote[%s] files err: %v", s.addr, addr, err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("%s getRemoteFiles status code:%v", s.addr, resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%s getRemoteFiles read body err:%v", s.addr, err)
		return
	}

	res := &LocalFilesRes{}
	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("%s getRemoteFiles body to json err:%v", s.addr, err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	serverFiles := map[string]serverFile{}
	// iterate the remote node file list
	for _, file := range res.Files {
		serverFile := serverFile{}
		serverFile.fileName = file
		serverFile.md5 = utils.MD5(file)
		serverFiles[file] = serverFile
	}

	s.files[host] = serverFiles
}

func (s *PServer) PING() string {
	return "PING"
}

func (s *PServer) PONG(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "PONG")
}

func (s *PServer) UploadFileTo(writer http.ResponseWriter, request *http.Request) {
	filename := request.URL.Query().Get("file")
	if filename == "" {
		//Get not set, send a 400 bad request
		http.Error(writer, "Get 'file' not specified in url.", 400)
		return
	}
	fmt.Println("Client requests: " + filename)

	//Check if file exists and open
	// Openfile, err := os.Open("files/" + Filename)
	Openfile, err := os.Open(s.filePath + "/" + filename)
	if err != nil {
		//File not found, send 404
		http.Error(writer, "File not found.", 404)
		return
	}
	defer Openfile.Close()
	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	writer.Header().Set("Content-Type", FileContentType)
	writer.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(writer, Openfile) //'Copy' the file to the client
}

func (s *PServer) SyncWithRemoteNode() {
	for _, localFile := range s.localFiles {
		flag := false
		for host, remoteFile := range s.files {
			if _, ok := remoteFile[localFile.fileName]; ok {
				flag = true
				break
				s.DownloadFileFrom(host, "/upload", localFile.fileName)
			}
		}

	}

	for host, remoteFile := range s.files {

		flag := false
		for _, localFile := range s.localFiles {
			if _, ok := remoteFile[localFile.fileName]; ok {
				flag = true
				break
			}
		}

		if !flag {
			s.DownloadFileFrom(host, "/upload", remoteFile.fileName)
		}

	}
}

func (s *PServer) SyncWithRemoteFileList() {
	for _, node := range s.nodes {
		s.getRemoteFiles(node, "/localFiles")
	}
}

// client for download file from remote server node
func (s *PServer) DownloadFileFrom(host, api, filename string) {
	addr := "http://" + host + api
	resp, err := http.Get(addr + "?file=" + filename)
	fmt.Printf("%s requests download file[%s] from %s", s.addr, filename, addr)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("request status code:", resp.StatusCode)
		return
	}

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(params)
	fileName := params["filename"]

	if s.isExistFile(fileName) {
		log.Println(s.addr + " have the same file name:" + filename)
		return
	}

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

	log.Printf("%s download file[%s] from node[%f] success:", s.addr, filename, host)
}

func (s *PServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/ping":
		s.PONG(w, req)
	case "/localFiles":
		s.GetLocalFileList(w, req)
	case "/upload":
		s.UploadFileTo(w, req)
	}
}

func (s *PServer) isExistFile(filename string) bool {
	flag := false
	md5Str := utils.MD5(filename)
	for _, file := range s.localFiles {
		if file.md5 == md5Str {
			flag = true
			break
		}
	}

	return flag
}

func Run(nfs NFSServer) (err error) {
	s := nfs.(*PServer)
	return http.ListenAndServe(s.addr, s)
}
