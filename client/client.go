package client

import (
	"fmt"
	"github.com/elaron/dmclock/request"
	"time"
)

type Client struct {
	name     string            `json:"name"`
	speed    int               `json:"speed"`
	wr       int               `json:"reservation"`
	wl       int               `json:"limit"`
	requests *request.Requests `json:"requests"`
}

func New(name string, speed, wr, wl int) *Client {
	return &Client{
		name:     name,
		speed:    speed,
		wr:       wr,
		wl:       wl,
		requests: request.New(),
	}
}

func (c *Client) RequestGenerator(reqCh chan string) {
	step := 1000 / c.speed
	fmt.Println(*c)
	ticker := time.NewTicker(time.Duration(step) * time.Millisecond)
	for range ticker.C {
		request.New()
		reqCh <- c.name
	}
}
