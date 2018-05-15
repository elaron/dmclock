package timer

/*
Timer structure is used to control the time behavior of client and server
*/
type Timer struct {
	//Ticker channel is used to send Ticker signal every 1 second
	Ticker chan bool
}

func New() *Timer {
	return &Timer{
		Ticker: make(chan bool, 1),
	}
}
