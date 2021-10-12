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
	"utils"
)

func (s *PServers) getRemoteFiles(host, api string) {
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

	serverFiles := map[string]PFile{}
	// iterate the remote node file list
	for _, file := range res.Files {
		serverFile := PFile{}
		serverFile.fileName = file
		serverFile.md5 = utils.MD5(file)
		serverFiles[file] = serverFile
	}

	s.localFiles = getPathFiles(s.filePath)
	s.files[host] = serverFiles
}

func (s *PServers) UploadFileTo(writer http.ResponseWriter, request *http.Request) {
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

// DownloadFileFrom client for download file from remote server node
func (s *PServers) DownloadFileFrom(host, api, filename string) {
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

	s.mu.Lock()
	defer s.mu.Unlock()
	downloadFile := PFile{}
	downloadFile.fileName = filename
	downloadFile.md5 = utils.MD5(filename)
	s.localFiles = append(s.localFiles, downloadFile)

	log.Printf("%s download file[%s] from node[%s] success:", s.addr, filename, host)
}

type LocalFilesRes struct {
	Files []string `json:"files"`
}

func (s *PServers) GetLocalFileList(w http.ResponseWriter, r *http.Request) {
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
}