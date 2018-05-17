package main

import (
	"github.com/elaron/dmclock/client"
	"github.com/elaron/dmclock/request"
	"github.com/elaron/dmclock/server"
	"time"
)

var g_clientList []*client.Client
var g_server *server.Server

func init() {
	g_clientList = []*client.Client{
		&client.Client{Name: "a", Speed: 800, Wr: 20, Wl: 70, Requests: []request.Req{}},
		&client.Client{Name: "b", Speed: 52, Wr: 30, Wl: 60, Requests: []request.Req{}},
		&client.Client{Name: "c", Speed: 40, Wr: 10, Wl: 70, Requests: []request.Req{}},
	}
	g_server = server.New(100)
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
