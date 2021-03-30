package main

import (
	"fmt"
	"pnfs"
	"time"
)

func main() {
	nodes := []string{"10.10.4.54:9998"}
	s := pnfs.New("10.10.4.54:9999", "./path", nodes)
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

	// s.PostLocalFileList("http://10.10.4.54:9998", "/getFileList", "./path")
	// pnfs.PostLocalFiles("10.10.4.54:9999", "")
	// s.ReceiveFileFrom()
	pnfs.Run(s)
}
