package main

import (
	"fmt"
	"pnfs"
)

func main() {
	nodes := []string{"10.10.4.54", "10.10.4.55"}
	s := pnfs.New("10.10.4.54:9999", "./path", nodes)
	if s != nil {
		fmt.Println("initialize success")
	}

	s.PostLocalFileList("http://10.10.4.54:9998", "/getFileList", "./path")
	// pnfs.PostLocalFiles("10.10.4.54:9999", "")
	// s.ReceiveFileFrom()
	pnfs.Run(s)
}
