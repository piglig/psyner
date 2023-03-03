package cmd

import "time"

type Client struct {
	relPathCheckSum map[string]string
	localDir        string
	tickerInterval  time.Duration
}
