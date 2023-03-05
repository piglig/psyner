package cmd

import (
	"fmt"
	"os"
	"psyner/client/taskrun"
	"time"
)

type Client struct {
	relPathCheckSum map[string]string
	config          ClientConfig
}

type ClientConfig struct {
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
	ticker := time.NewTicker(c.config.TickerInterval)
	defer ticker.Stop()

	go taskrun.CheckLocalDirChecksum(c.config.LocalDir, ticker)

}