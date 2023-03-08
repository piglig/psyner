package cmd

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"psyner/client/taskrun"
	"strings"
	"time"
)

type Client struct {
	relPathCheckSum map[string]string
	config          ClientConfig
}

type ClientConfig struct {
	ServerAddr     string
	LocalDir       string
	TickerInterval time.Duration
}

func NewClient(config ClientConfig) (*Client, error) {
	if config.LocalDir == "" {
		return nil, fmt.Errorf("local dir %s not invalid", config.LocalDir)
	}

	_, err := os.Stat(config.LocalDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("local dir %s not exist", config.LocalDir)
		} else {
			return nil, fmt.Errorf("local dir stat invalid %v", err)
		}
	}

	if config.TickerInterval == 0 {
		config.TickerInterval = 5 * time.Second
	}

	return &Client{
		relPathCheckSum: make(map[string]string),
		config:          config,
	}, err
}

func (c *Client) Start() {
	// connect to server
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go taskrun.CheckLocalDirChecksum(c.config.LocalDir, c.config.TickerInterval)

	for {
		// read available files from server
		var fileName string
		err = gob.NewDecoder(conn).Decode(&fileName)
		if err != nil {
			log.Fatal(err)
		}
		fileName = strings.TrimSpace(fileName)

		// send selected file name to server
		err = gob.NewEncoder(conn).Encode(fileName)
		if err != nil {
			log.Fatal(err)
		}

		// receive file data from server
		fileData := bytes.Buffer{}
		_, err = io.Copy(&fileData, conn)
		if err != nil {
			log.Fatal(err)
		}

		// save file to local computer
		err = os.WriteFile(fileName, fileData.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Saved file %s\n", fileName)
	}
}
