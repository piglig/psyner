package cli

import (
	"flag"
	"net"
	"os"
	"path/filepath"
)

const (
	MasterFlag = "m" // master node
	HostFlag   = "h" // host
	PortFlag   = "p" // port
	DirFlag    = "d" // data directory
)

type PNFSFlag struct {
	master bool
	host   string
	port   int
	dir    string
}

func InitFlag() *PNFSFlag {
	pnfsFlag := new(PNFSFlag)
	flag.BoolVar(&pnfsFlag.master, MasterFlag, false, "master node")
	flag.StringVar(&pnfsFlag.host, HostFlag, "127.0.0.1", "pnfs host")
	flag.IntVar(&pnfsFlag.port, PortFlag, 3100, "pnfs port")
	flag.StringVar(&pnfsFlag.dir, DirFlag, "/path", "pnfs file directory")
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

	return true
}
