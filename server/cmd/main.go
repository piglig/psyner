package main

import (
	"log"
	"psyner/server/ctx"
	_ "psyner/server/taskrun/action/event"
)

const (
	listenAddr = "localhost:8888"
	dir        = "./data"
)

func main() {
	server, err := ctx.NewServer(ctx.ServerConfig{
		// host.docker.internal
		ListenAddr: listenAddr,
		LocalDir:   dir,
	})
	if err != nil {
		log.Fatal("NewServer", err.Error())
	}

	server.Run()
}
