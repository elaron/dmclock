package request

import "sync"

type Req struct {
	ClientName string
	ReqId      int
}

type Requests struct {
	list []Req
}

func New() *Requests  {
	return &Requests{
		list: []Req{},
	}
}

func (r *Requests)AddReq(clientName string) int {
	r.list = append(r.list, re)
}

var g_reqId int
var g_reqIdLock sync.RWMutex

func NewId() int {
	g_reqIdLock.Lock()
	g_reqId += 1
	g_reqIdLock.Unlock()
	return g_reqId
}
