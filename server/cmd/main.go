package main

import (
	"log"
	"psyner/server/taskrun/runner"
)

const (
	listenAddr = "localhost:8888"
	dir        = "./data"
)

func main() {
	server, err := runner.NewServer(runner.ServerConfig{
		// host.docker.internal
		ListenAddr: listenAddr,
		LocalDir:   dir,
	})
	if err != nil {
		log.Fatal("NewServer", err.Error())
	}

	server.Run()
}
