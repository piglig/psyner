package cli

import (
	"flag"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	MasterFlag = "m" // master node
	HostFlag   = "h" // host
	PortFlag   = "p" // port
	DirFlag    = "d" // data directory
	SlaveFlag  = "s" // slave node
)

type PNFSFlag struct {
	master  bool
	slaveOf string
	host    string
	port    int
	dir     string
}

func InitFlag() *PNFSFlag {
	pnfsFlag := new(PNFSFlag)
	flag.BoolVar(&pnfsFlag.master, MasterFlag, false, "master node")
	flag.StringVar(&pnfsFlag.slaveOf, SlaveFlag, "", "slave of master node addr, format: host:port")
	flag.StringVar(&pnfsFlag.host, HostFlag, "127.0.0.1", "server host")
	flag.IntVar(&pnfsFlag.port, PortFlag, 3100, "server port")
	flag.StringVar(&pnfsFlag.dir, DirFlag, "/path", "server file directory")
	flag.Parse()

	if !pnfsFlag.checkFlag() {
		return nil
	}

	return pnfsFlag
}

func (p *PNFSFlag) checkFlag() bool {
	// host
	if net.ParseIP(p.host) == nil {
		flag.Usage()
		return false
	}

	p.dir = filepath.Join(".", p.dir)
	// path
	dir, err := os.Stat(p.dir)
	if err != nil || !dir.IsDir() {
		flag.Usage()
		return false
	}

	// check master status
	if len(p.slaveOf) > 0 {
		conn, err := net.DialTimeout("tcp", p.slaveOf, time.Second*3)
		if err != nil {
			flag.Usage()
			return false
		}

		if conn != nil {
			defer conn.Close()
		}
	}

	return true
}

func (p *PNFSFlag) IsMaster() bool {
	return p.master
}

func (p *PNFSFlag) GetFilePath() string {
	return p.dir
}

func (p *PNFSFlag) GetHostPort() (string, int) {
	return p.host, p.port
}

func (p *PNFSFlag) GetMasterAddr() string {
	return p.slaveOf
}
