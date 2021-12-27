package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"pnfs/cli"
	"sync"
	"pnfs/utils"
)

const (
	UploadAPI     = "/upload"
	LocalFilesAPI = "/files"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type NfsServerFunc interface {
	// UploadFileTo server for other server node download
	UploadFileTo(writer http.ResponseWriter, request *http.Request)
	GetLocalFileList(w http.ResponseWriter, r *http.Request)
	// ServeHTTP use for http server
	ServeHTTP(w http.ResponseWriter, r *http.Request)

	// Ping check server is normal
	Ping(w http.ResponseWriter, r *http.Request)
	// Pong response server check status
	Pong(w http.ResponseWriter, r *http.Request)
}

type PNfs struct {
	servers      []*PServer
	rwm          sync.RWMutex
	fileToServer FileToServer
	masterAddr   string
}

type FileToServer map[string]PFile

// New initial server server
func New(flag cli.PNFSFlag) (*PNfs, error) {
	res := &PNfs{}
	if flag.IsMaster() {
		res.fileToServer = make(FileToServer, 0)

		return res, nil
	}

	// master addr
	res.masterAddr = flag.GetMasterAddr()
	host, port := flag.GetHostPort()
	server := PServer{
		host:     host,
		port:     port,
		active:   false,
		fsPath:   flag.GetFilePath(),
		isMaster: false,
	}



	for
}


type PServer struct {
	host     string
	port     int
	active   bool
	fsPath   string
	isMaster bool
	files    []*PFile
}

func (p *PServer) IsActive() bool {
	return p.active
}


type PFile struct {
	file os.FileInfo
	md5  string
}

func GetPFileFromDir(dir string) ([]*PFile, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	res := make([]*PFile, 10)
	for _, file := range files {
		p := &PFile{}
		p.file = file
		p.md5 = utils.MD5(file.Name())

		res = append(res, p)
	}
	return res, nil
}

// PFile server file struct
//type PFile struct {
//	FileName string
//	FileInfo os.FileInfo
//	Md5      string
//
//	// ServerIndex the file locate at server id
//	ServerIndex int
//}

//type PServer struct {
//	// LocalFiles []PFile // the server node files
//	Alive bool     // the server alive status
//	Url   *url.URL // the server addr
//}

type Temp struct {
	servers []*PServer
	files   []*PFile

	mu sync.Mutex // protects currently request
}

type PServers struct {
	servers []*PServer

	files      map[string]map[string]PFile // the other server node and files md5 string
	nodes      []string                    // the server node host
	filePath   string                      // the server node file path
	localFiles []PFile                     // the current server node files

	addr string // the server node host and port

	rwLock sync.RWMutex // rw lock
	mu     sync.Mutex   // protects currently request
}



/*func New(addr, path string, nodes []*url.URL) *PServers {
	s := &PServers{
		addr:     addr,
		filePath: path,
	}
	for _, node := range nodes {
		server := &PServer{
			Alive: true,
			Url:   node,
		}

		s.servers = append(s.servers, server)
	}
	fmt.Printf("addr [%s], local file path[%s], server nodes%v\n", s.addr, s.filePath, nodes)
	return s
}*/

func (p *PServer) GetLocalFileList() {

}

func getPathFiles(filePath string) []PFile {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		log.Fatalf("read path[%s] files error: %v", filePath, err)
		return nil
	}

	serverFiles := []PFile{}
	for _, file := range files {
		serverFile := PFile{}
		// serverFile.fileInfo = file
		serverFile.FileName = file.Name()
		serverFile.Md5 = utils.MD5(file.Name())
		serverFiles = append(serverFiles, serverFile)
	}
	return serverFiles
}

func (s *PServers) HealthCheck() {
	addr := "http://"
	resp, err := http.Get(addr)

	if err != nil {
		index := utils.GetAddrIndexFromNodes(host, s.nodes)
		log.Printf("%s remove disconnect node:%s\n", s.addr, s.nodes[index])
		s.nodes = append(s.nodes[:index], s.nodes[index+1:]...)
		// log.Printf("%s ping remote[%s] err: %v", s.addr, addr, err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("%s ping remote status code:%v\n", s.addr, resp.StatusCode)
		return
	}

	if !utils.IsAddrInNodes(host, s.nodes) {
		s.nodes = append(s.nodes, host)
		log.Printf("%s add new node:%s\n", s.addr, host)
		return
	}
}

func (s *PServers) SyncWithRemoteNode() {
	for host, remoteFile := range s.files {
		for fileName := range remoteFile {
			flag := false
			for _, localFile := range s.localFiles {
				if localFile.FileName == fileName {
					flag = true
					break
				}
			}

			if !flag {
				s.DownloadFileFrom(host, UploadAPI, fileName)
			}
		}
	}
}

func (s *PServers) SyncWithRemoteFileList() {
	for _, node := range s.nodes {
		s.getRemoteFiles(node, LocalFilesAPI)
	}
}

func (s *PServers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case LocalFilesAPI:
		s.GetLocalFileList(w, req)
	case UploadAPI:
		s.UploadFileTo(w, req)
	}
}

func (s *PServers) isExistFile(filename string) bool {
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

func Run(nfs NFSServerFunc) (err error) {
	s := nfs.(*PServers)
	return http.ListenAndServe(s.addr, s)
}
