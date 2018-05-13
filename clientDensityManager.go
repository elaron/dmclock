package main

import (
	"sort"
	"sync"
)

type ClientName string

type ClientDensity struct {
	name        ClientName
	Wr          int
	Wl          int
	currDensity int
}

type ClientDensityManager struct {
	m    map[ClientName]*ClientDensity
	l    []*ClientDensity
	lock sync.RWMutex
	less []lessFunc
}

func NewClientDensityManage() *ClientDensityManager {
	return &ClientDensityManager{
		m: make(map[ClientName]*ClientDensity),
		l: []*ClientDensity{},
		less: []lessFunc{
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
			},
		},
	}
}

func (cdm *ClientDensityManager) UpdateClientDensity(name ClientName, wr, wl, num int) {
	cdm.lock.Lock()
	cd, ok := cdm.m[name]
	if !ok {
		clientDensity := &ClientDensity{name: name, Wr: wr, Wl: wl, currDensity: num}
		cdm.m[name] = clientDensity
		cdm.l = append(cdm.l, clientDensity)
	} else {
		cd.currDensity += num
	}
	cdm.lock.Unlock()
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
