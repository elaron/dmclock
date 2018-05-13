package main

import (
	"fmt"
	"sync"
	"time"
)

type Req struct {
	ClientName string
	ReqId      int
}

type Client struct {
	Name     string `json:"name"`
	Speed    int    `json:"speed"`
	Wr       int    `json:"reservation"`
	Wl       int    `json:"limit"`
	Requests []Req  `json:"requests"`
}

var g_reqId int
var g_reqIdLock sync.RWMutex

func NewId() int {
	g_reqIdLock.Lock()
	g_reqId += 1
	g_reqIdLock.Unlock()
	return g_reqId
}

func (c *Client) RequestGenerator(reqCh chan string) {
	step := 1000 / c.Speed
	fmt.Println(*c)
	ticker := time.NewTicker(time.Duration(step) * time.Millisecond)
	for range ticker.C {
		reqCh <- c.Name
	}
}

type Server struct {
	Capacity      int   `json:"capacity"`
	WaitQueue     []Req `json:"wait_queue"`
	waitQueueLock sync.RWMutex
}

func (s *Server) Enqueue(reqCh chan string) {
	for {
		select {
		case client := <-reqCh:
			//s.waitQueueLock.Lock()
			s.WaitQueue = append(s.WaitQueue, Req{client, NewId()})
			//s.waitQueueLock.Unlock()
		}
	}
}

func (s *Server) FIFODequeue() {
	counter := make(map[string]int)
	max := s.Capacity
	//s.waitQueueLock.RLock()
	for _, v := range s.WaitQueue {
		if max <= 0 {
			break
		}
		//fmt.Println(v)
		n, ok := counter[v.ClientName]
		if !ok {
			counter[v.ClientName] = 1
		} else {
			counter[v.ClientName] = n + 1
		}
		max -= 1
	}
	//s.waitQueueLock.RUnlock()
	fmt.Println(len(s.WaitQueue), counter)
}

var g_clientList []*Client
var g_server Server

func init() {
	g_reqId = 0
	g_clientList = []*Client{
		&Client{Name: "a", Speed: 800, Wr: 20, Wl: 70, Requests: []Req{}},
		&Client{Name: "b", Speed: 52, Wr: 30, Wl: 60, Requests: []Req{}},
		&Client{Name: "c", Speed: 40, Wr: 10, Wl: 70, Requests: []Req{}},
	}
	g_server = Server{Capacity: 100, WaitQueue: []Req{}}
}

func main() {
	reqCh := make(chan string, 100)
	go g_server.Enqueue(reqCh)
	for _, client := range g_clientList {
		go client.RequestGenerator(reqCh)
	}

	time.Sleep(1 * time.Second)
	g_server.FIFODequeue()
}
