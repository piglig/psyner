package main

import (
	"log"
	"psyner/server/cmd"
	_ "psyner/server/taskrun/action/event"
)

const (
	listenAddr = "localhost:8888"
	dir        = "./data"
)

func main() {
	server, err := cmd.NewServer(cmd.ServerConfig{
		// host.docker.internal
		ListenAddr: listenAddr,
		LocalDir:   dir,
	})
	if err != nil {
		log.Fatal("NewServer", err.Error())
	}

	server.Run()
}
