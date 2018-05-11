package main

type ReqId int
type Client struct {
	Name     string  `json:"name"`
	Wr       int     `json:"reservation"`
	Wl       int     `json:"limit"`
	Requests []ReqId `json:"requests"`
}

type Server struct {
	Capacity  int     `json:"capacity"`
	WaitQueue []ReqId `json:"wait_queue"`
}

var g_clientList []Client
var g_server Server

func init() {
	g_clientList = []Client{
		Client{"a", 20, 70, []ReqId{}},
		Client{"b", 30, 60, []ReqId{}},
		Client{"c", 10, 70, []ReqId{}},
	}
	g_server = Server{100, []ReqId{}}
}

func main() {

}
