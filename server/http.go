package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"pnfs/utils"
	"strconv"
)

const (
	FileListAPI = "files"
)

type Data struct {
	Code int
	Msg  string
	Data interface{}
}

func result(w http.ResponseWriter, code int, msg string, data interface{}) {
	res := Data{Code: code, Msg: msg, Data: data}
	dataBytes, _ := json.Marshal(&res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dataBytes)
}

func (p *PServer) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG")
}

func (p *PServer) GetUnSyncedFileFromMaster(w http.ResponseWriter, r *http.Request) {
	if p.isMaster {
	} else {
		fileURL := p.masterAddr + "/" + FileListAPI
		resp, err := http.Get(fileURL)
		if err != nil {
			log.Println("GetUnSyncedFile", "err", err)
			result(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			result(w, http.StatusInternalServerError, err.Error(), nil)
			log.Fatal(err)
		}
		f := PFile{}
		if err = json.Unmarshal(body, &f); err != nil {
			log.Println("GetUnSyncedFile", "err", err)
			result(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		exist := false
		for _, localFile := range p.files {
			if localFile.md5 == f.md5 {
				exist = true
				break
			}
		}

		if !exist {
			p.rwm.Lock()
			defer p.rwm.Unlock()
			p.files = append(p.files, f)
		}
	}
	result(w, http.StatusOK, http.StatusText(http.StatusOK), nil)
}

func GetHostPort(r *http.Request) (string, string) {
	host := r.URL.Query().Get("host")
	port := r.URL.Query().Get("port")
	return host, port
}

func (p *PNfs) GetServer(host, port string) PServer {
	for _, server := range p.servers {
		if net.JoinHostPort(server.host, server.port) == net.JoinHostPort(host, port) {
			return server
		}
	}
	return PServer{host: host, port: port}
}

func (p *PNfs) GetUnSyncedFile(w http.ResponseWriter, r *http.Request) {
	host, port := GetHostPort(r)

	exist := false
	for _, s := range p.servers {
		if s.host == host && s.port == port {
			exist = true
			break
		}
	}

	if exist {
		file := PFile{}
		p.rwm.RLock()
		defer p.rwm.RUnlock()
		if len(p.files) == 0 {
			result(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
			return
		}

		serverFiles := p.serverToFiles[net.JoinHostPort(host, port)]

		for _, f := range p.files {
			if !IsExistInFiles(f, serverFiles) {
				file = f
				res, _ := json.Marshal(file)
				result(w, http.StatusOK, http.StatusText(http.StatusOK), res)
				return
			}
		}
		result(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
		return
	} else {
		server := p.GetServer(host, port)
		p.addServer(server)
		result(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
		return
	}
}

func (p *PNfs) IsExistInFiles(f PFile) bool {
	for _, file := range p.files {
		if file.md5 == f.md5 && f.file.Name() == file.file.Name() {
			return true
		}
	}
	return false
}

func IsExistInFiles(file PFile, files []PFile) bool {
	for _, f := range files {
		if file.md5 == f.md5 && f.file.Name() == file.file.Name() {
			return true
		}
	}
	return false
}

func (p *PNfs) addServer(server PServer) {
	exist := false
	for _, s := range p.servers {
		if s.host == server.host && s.port == server.port {
			exist = true
			break
		}
	}

	if !exist {
		p.rwm.Lock()
		defer p.rwm.Unlock()
		p.servers = append(p.servers, server)
	}
}

func (s *main.PServers) getRemoteFiles(host, api string) {
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

	serverFiles := map[string]main.PFile{}
	// iterate the remote node file list
	for _, file := range res.Files {
		serverFile := main.PFile{}
		serverFile.FileName = file
		serverFile.Md5 = utils.MD5(file)
		serverFiles[file] = serverFile
	}

	s.localFiles = main.getPathFiles(s.filePath)
	s.files[host] = serverFiles
}

func (s *main.PServers) UploadFileTo(writer http.ResponseWriter, request *http.Request) {
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
func (s *main.PServers) DownloadFileFrom(host, api, filename string) {
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
	downloadFile := main.PFile{}
	downloadFile.FileName = filename
	downloadFile.Md5 = utils.MD5(filename)
	s.localFiles = append(s.localFiles, downloadFile)

	log.Printf("%s download file[%s] from node[%s] success:", s.addr, filename, host)
}

type LocalFilesRes struct {
	Files []string `json:"files"`
}

func (s *main.PServers) GetLocalFileList(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := &LocalFilesRes{}
	for _, file := range s.localFiles {
		res.Files = append(res.Files, file.FileName)
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Printf("postLocalFiles marshal to json err:%v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)
}
