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
	client := &Client{
		name:     name,
		speed:    speed,
		wr:       wr,
		wl:       wl,
		requests: request.New(),
	}
	go client.RequestGenerator()
	return client
}

func (c *Client) RequestGenerator() {
	step := 1000 / c.speed
	fmt.Println(*c)
	ticker := time.NewTicker(time.Duration(step) * time.Millisecond)
	for range ticker.C {
		c.requests.AddReq(c.name)
	}
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) Wr() int {
	return c.wr
}

func (c *Client) Wl() int {
	return c.wl
}

func (c *Client) Speed() int {
	return c.speed
}
func (c *Client) RequestsNum() int {
	return c.requests.Len()
}

func (c *Client) DoOneRequest() *request.Req {
	return c.requests.Pop()
}
