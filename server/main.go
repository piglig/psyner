package main

import (
	"log"
	"psyner/server/taskrun"
)

const (
	listenAddr = "localhost:8888"
	dir        = "./data"
)

func main() {
	server, err := taskrun.NewServer(taskrun.ServerConfig{
		// host.docker.internal
		ListenAddr: listenAddr,
		LocalDir:   dir,
	})
	if err != nil {
		log.Fatal("NewServer", err.Error())
	}

	server.Run()
}
