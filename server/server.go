package server

import (
	r "dmclock/request"
	"fmt"
	"sync"
)

type Server struct {
	Capacity      int     `json:"capacity"`
	WaitQueue     []r.Req `json:"wait_queue"`
	waitQueueLock sync.RWMutex
	cdm           *ClientDensityManager
}

func New(cap int) *Server {
	return &Server{
		Capacity:  cap,
		WaitQueue: []r.Req{},
		cdm:       NewClientDensityManage(),
	}
}

func (s *Server) Enqueue(reqCh chan string) {
	for {
		select {
		case client := <-reqCh:
			//s.waitQueueLock.Lock()
			s.WaitQueue = append(s.WaitQueue, r.Req{client, r.NewId()})
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
