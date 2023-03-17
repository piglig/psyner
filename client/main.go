package main

import (
	"log"
	"psyner/client/cmd"
	"time"
)

func main() {
	client, err := cmd.NewClient(cmd.ClientConfig{
		// host.docker.internal
		ServerAddr:     ":8888",
		LocalDir:       "./client/data",
		TickerInterval: 10 * time.Second,
	})
	if err != nil {
		log.Fatal("NewClient", err.Error())
	}

	client.Start()
}
