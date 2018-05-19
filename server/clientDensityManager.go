package server

import (
	"errors"
	"github.com/elaron/dmclock/client"
	"github.com/elaron/dmclock/request"
	"github.com/elaron/dmclock/timer"
	"log"
	"sort"
	"sync"
)

var (
	ClientExistErr    error = errors.New("Client already exist")
	ClientNotExistErr error = errors.New("Client not exist")
)

type ClientName string

type ClientDensity struct {
	*client.Client
	currDensity int
}

type ClientDensityManager struct {
	m    map[ClientName]*ClientDensity
	l    []*ClientDensity
	lock sync.RWMutex
	less []lessFunc
	t    *timer.Timer
}

func (cd *ClientDensity) ReservationRate() float32 {
	return float32(cd.currDensity) / float32(cd.Wr())
}

func (cd *ClientDensity) LimitationRate() float32 {
	return float32(cd.currDensity) / float32(cd.Wl())
}

func NewClientDensityManage() *ClientDensityManager {
	cdm := &ClientDensityManager{
		m: make(map[ClientName]*ClientDensity),
		l: []*ClientDensity{},
		t: timer.New(),
		less: []lessFunc{
			func(c1 *ClientDensity, c2 *ClientDensity) bool {
				len1 := c1.RequestsNum()
				len2 := c2.RequestsNum()
				if len1 > 0 && len2 <= 0 {
					return true
				} else {
					return false
				}
			},
			func(c1 *ClientDensity, c2 *ClientDensity) bool {
				rateRc1, rateRc2 := c1.ReservationRate(), c2.ReservationRate()
				if rateRc1 >= 1 && rateRc2 >= 1 {
					return false
				} else {
					return rateRc1 < rateRc2
				}
			},
			func(c1 *ClientDensity, c2 *ClientDensity) bool {
				rateLc1, rateLc2 := c1.LimitationRate(), c2.LimitationRate()
				return rateLc1 < rateLc2
			},
		},
	}

	go func() {
		for {
			<-cdm.t.Ticker
			cdm.ResetClientDensity()
		}
	}()

	return cdm
}

func (cdm *ClientDensityManager) GetClient(name ClientName) (*ClientDensity, bool) {
	cdm.lock.RLock()
	c, ok := cdm.m[name]
	cdm.lock.RUnlock()
	return c, ok
}

func (cdm *ClientDensityManager) PutClient(name ClientName, cd *ClientDensity) {
	cdm.lock.Lock()
	cdm.m[name] = cd
	cdm.l = append(cdm.l, cd)
	cdm.lock.Unlock()
}

func (cdm *ClientDensityManager) DeleteClient(name ClientName) {
	cdm.lock.Lock()
	delete(cdm.m, name)
	idx := 0
	for i := 0; i < len(cdm.l); i++ {
		if cdm.l[i] != nil && cdm.l[i].Name() == string(name) {
			idx = i
			break
		}
	}
	if idx < len(cdm.l) {
		cdm.l = append(cdm.l[0:idx], cdm.l[idx+1:]...)
	} else {
		log.Printf("Never find %s in ClientDensityManager\n", name)
	}
	cdm.lock.Unlock()
}

func (cdm *ClientDensityManager) AddClient(name ClientName, client *client.Client) error {

	_, ok := cdm.GetClient(name)
	if true == ok {
		log.Println(name, "client already exsit.")
		return ClientExistErr
	}

	cd := &ClientDensity{client, 0}
	cdm.PutClient(name, cd)

	return nil
}

func (cdm *ClientDensityManager) scheduleClient() *request.Req {
	cdm.Sort()

	//select a client
	if len(cdm.l) == 0 {
		log.Println("Client list is empty")
		return nil
	}
	client := cdm.l[0]
	client.currDensity += 1

	//pop its first request
	return client.DoOneRequest()
}

func (cdm *ClientDensityManager) UpdateClientDensity(name ClientName, density int) error {
	client, ok := cdm.GetClient(name)
	if false == ok {
		log.Printf("client %s not exist\n", name)
		return ClientNotExistErr
	}

	client.currDensity = density
	return nil
}

func (cdm *ClientDensityManager) ResetClientDensity() {
	cdm.lock.Lock()
	for _, cd := range cdm.l {
		cd.currDensity = 0
	}
	cdm.lock.Unlock()
}

type lessFunc func(c1 *ClientDensity, c2 *ClientDensity) bool

func (c *ClientDensityManager) Sort() {
	sort.Sort(c)
	c.PrintClients()
}

func (cdm *ClientDensityManager) PrintClients() {
	for i, client := range cdm.l {
		log.Printf("idx=%d %s reqLen=%d rateR=%f rateL=%f currentDensity=%d\n",
			i, client.Name(), client.RequestsNum(), client.ReservationRate(), client.LimitationRate(), client.currDensity)
	}
}

func (c *ClientDensityManager) Len() int {
	return len(c.l)
}

func (c *ClientDensityManager) Swap(i, j int) {
	c.l[j], c.l[i] = c.l[i], c.l[j]
}
func (c *ClientDensityManager) Less(i, j int) bool {
	p, q := c.l[i], c.l[j]
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
