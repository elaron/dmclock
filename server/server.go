package server

import (
	"fmt"
	"github.com/elaron/dmclock/client"
	r "github.com/elaron/dmclock/request"
	"log"
	"sync"
)

type Server struct {
	capacity      int     `json:"capacity"`
	WaitQueue     []r.Req `json:"wait_queue"`
	waitQueueLock sync.RWMutex
	cdm           *ClientDensityManager
}

func New(cap int) *Server {
	return &Server{
		capacity:  cap,
		WaitQueue: []r.Req{},
		cdm:       NewClientDensityManage(),
	}
}

func (s *Server) AddClient(name ClientName, client *client.Client) {
	s.cdm.AddClient(name, client)
}

func (s *Server) FIFODequeue() {
	counter := make(map[string]int)
	max := s.capacity
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

func (s *Server) DensityDequeue() {
	counter := make(map[string]int)
	for i := 0; i < s.capacity; i++ {
		req := s.cdm.scheduleClient()
		if nil == req {
			continue
		}
		counter[req.ClientName] += 1
		log.Println(req.ClientName, req.ReqId)
	}
	fmt.Println(counter)
}
