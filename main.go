package main

import (
	"fmt"
	"sort"
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

type ClientDensity struct {
	Wr          int
	Wl          int
	currDensity int
}

type lessFunc func(c1 *ClientDensity, c2 *ClientDensity) bool

type SortClient struct {
	cDensity []ClientDensity
	less     []lessFunc
}

func (c *SortClient) Sort(cDen []ClientDensity) {
	c.cDensity = cDen
	sort.Sort(c)
}

func NewClientSorter(less ...lessFunc) *SortClient {
	return &SortClient{
		less: less,
	}
}
func (c *SortClient) Len() int {
	return len(c.cDensity)
}

func (c *SortClient) Swap(i, j int) {
	c.cDensity[j], c.cDensity[i] = c.cDensity[i], c.cDensity[j]
}
func (c *SortClient) Less(i, j int) bool {
	p, q := &c.cDensity[i], &c.cDensity[j]
	for _, less := range c.less {
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		default:
			continue
		}
	}
	return false
}

type Server struct {
	Capacity      int   `json:"capacity"`
	WaitQueue     []Req `json:"wait_queue"`
	waitQueueLock sync.RWMutex
	clients       []ClientDensity
	ClientSorter  *SortClient
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

func DensityDequeue() {

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

	g_server.ClientSorter = NewClientSorter(
		func(c1 *ClientDensity, c2 *ClientDensity) bool {
			rateRc1, rateRc2 := float32(c1.currDensity)/float32(c1.Wr), float32(c2.currDensity)/float32(c2.Wr)
			if rateRc1 >= 1 && rateRc2 >= 1 {
				return false
			} else {
				return rateRc1 < rateRc2
			}
		},
		func(c1 *ClientDensity, c2 *ClientDensity) bool {
			rateLc1, rateLc2 := float32(c1.currDensity)/float32(c1.Wl), float32(c2.currDensity)/float32(c2.Wl)
			return rateLc1 < rateLc2
		})
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
