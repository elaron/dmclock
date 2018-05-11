package main

import (
	"fmt"
	"sync"
	"time"
)

type ReqId string
type Client struct {
	Name     string  `json:"name"`
	Speed    int     `json:"speed"`
	Wr       int     `json:"reservation"`
	Wl       int     `json:"limit"`
	Requests []ReqId `json:"requests"`
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
	for {
		step := 1000 / c.Speed
		ticker := time.NewTicker(time.Duration(step) * time.Microsecond)
		for range ticker.C {
			reqCh <- c.Name
		}
	}
}

type Server struct {
	Capacity      int     `json:"capacity"`
	WaitQueue     []ReqId `json:"wait_queue"`
	waitQueueLock sync.RWMutex
}

func (s *Server) Enqueue(reqCh chan string) {
	for {
		select {
		case client := <-reqCh:
			s.waitQueueLock.Lock()
			s.WaitQueue = append(s.WaitQueue, ReqId(fmt.Sprintf("%s.%d", client, NewId())))
			s.waitQueueLock.Unlock()
		}
	}
}

func (s *Server) PrintQueue() {
	s.waitQueueLock.RLock()
	for _, v := range s.WaitQueue {
		fmt.Println(v)
	}
	s.waitQueueLock.RUnlock()
}

var g_clientList []*Client
var g_server Server

func init() {
	g_reqId = 0
	g_clientList = []*Client{
		&Client{Name: "a", Speed: 50, Wr: 20, Wl: 70, Requests: []ReqId{}},
		&Client{Name: "b", Speed: 52, Wr: 30, Wl: 60, Requests: []ReqId{}},
		&Client{Name: "c", Speed: 40, Wr: 10, Wl: 70, Requests: []ReqId{}},
	}
	g_server = Server{Capacity: 100, WaitQueue: []ReqId{}}
}

func main() {
	reqCh := make(chan string, 100)
	for _, client := range g_clientList {
		go client.RequestGenerator(reqCh)
	}
	go g_server.Enqueue(reqCh)

	time.Sleep(1 * time.Second)
	g_server.PrintQueue()
}
