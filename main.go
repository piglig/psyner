package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"pnfs"
	"strings"
	"time"
)

func main() {

	var nodeList string
	var port int
	var path string
	flag.StringVar(&nodeList, "nodes", "", "PNFS other nodes, use commas to separate")
	flag.StringVar(&path, "path", "./path", "PNFS local path for synchronize")
	flag.IntVar(&port, "port", 3100, "Port to serve")
	flag.Parse()

	if len(nodeList) == 0 {
		log.Fatal("Please provide one or more nodes to synchronize")
	}

	var servers []*url.URL
	nodes := strings.Split(nodeList, ",")
	for _, node := range nodes {
		serverUrl, err := url.Parse(node)
		if err != nil {
			log.Fatal(err)
		}
		servers = append(servers, serverUrl)
	}

	nodes := []string{"10.10.4.54:9998"}
	s := pnfs.New("10.10.4.54:9999", "./path", servers)
	if s != nil {
		fmt.Println("initialize success")
	}

	ticker := time.NewTicker(3 * time.Second)
	ticker2 := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker2.C:
				s.SyncWithRemoteFileList()
				ticker2 = time.NewTicker(1 * time.Second)
			case <-ticker.C:
				s.SyncWithRemoteNode()
				ticker = time.NewTicker(3 * time.Second)
			case <-quit:
				ticker.Stop()
				return
			}
		}

	}()

	pnfs.Run(s)
}
