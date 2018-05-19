package main

import (
	"github.com/elaron/dmclock/client"
	"github.com/elaron/dmclock/server"
	"time"
)

var g_server *server.Server

func init() {
	g_server = server.New(100)
}

func main() {
	g_server.AddClient("a", client.New("a", 800, 20, 70))
	g_server.AddClient("b", client.New("b", 50, 30, 60))
	g_server.AddClient("c", client.New("c", 40, 10, 70))

	time.Sleep(1 * time.Second)
	g_server.DensityDequeue()
	//g_server.FIFODequeue()
}
