package client

import (
	"dmclock/request"
	"fmt"
	"time"
)

type Client struct {
	Name     string        `json:"name"`
	Speed    int           `json:"speed"`
	Wr       int           `json:"reservation"`
	Wl       int           `json:"limit"`
	Requests []request.Req `json:"requests"`
}

func (c *Client) RequestGenerator(reqCh chan string) {
	step := 1000 / c.Speed
	fmt.Println(*c)
	ticker := time.NewTicker(time.Duration(step) * time.Millisecond)
	for range ticker.C {
		reqCh <- c.Name
	}
}
