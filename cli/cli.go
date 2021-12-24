package cli

import (
	"flag"
	"net"
	"os"
)

const (
	MasterFlag = "m" // master node
	HostFlag   = "h" // host
	PortFlag   = "p" // port
	DirFlag    = "d" // data directory
)

type masterFlag struct {
	set bool
}

func (m *masterFlag) Set(value string) error {
	m.set = true
	return nil
}

func (m *masterFlag) String() string {
	return ""
}

type PNFSFlag struct {
	master masterFlag
	host   string
	port   int
	dir    string
}

func InitFlag() *PNFSFlag {
	pnfsFlag := &PNFSFlag{}
	flag.Var(&pnfsFlag.master, MasterFlag, "master node")
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
	// master node
	if !p.master.set {
		flag.Usage()
		return false
	}

	// host
	if net.ParseIP(p.host) == nil {
		flag.Usage()
		return false
	}

	// path
	dir, err := os.Stat(p.dir)
	if err != nil || !dir.IsDir() {
		flag.Usage()
		return false
	}

	return true
}
