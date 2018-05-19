package request

import (
	"log"
	"sync"
)

type Req struct {
	ClientName string
	ReqId      int
}

type Requests struct {
	list []Req
	lock sync.RWMutex
}

func New() *Requests {
	return &Requests{
		list: []Req{},
	}
}

func (r *Requests) AddReq(clientName string) int {
	id := NewId()
	r.lock.Lock()
	r.list = append(r.list, Req{ClientName: clientName, ReqId: id})
	r.lock.Unlock()
	return id
}

func (r *Requests) Pop() *Req {
	if len(r.list) == 0 {
		log.Println("Request list is empty")
		return nil
	}
	r.lock.Lock()
	req := r.list[0]
	r.list = r.list[1:]
	r.lock.Unlock()
	return &req
}

var g_reqId int
var g_reqIdLock sync.RWMutex

func NewId() int {
	g_reqIdLock.Lock()
	g_reqId += 1
	g_reqIdLock.Unlock()
	return g_reqId
}

func (r *Requests) Len() int {
	return len(r.list)
}
