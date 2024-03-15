package main

import (
	"os"
	"playful-patterns.com/bakoko/gui"
	. "playful-patterns.com/bakoko/networking"
	"time"
)

func main() {
	var worldProxyTcpIp WorldProxyTcpIp
	worldProxyTcpIp.Endpoint = os.Args[1] // localhost:56901 or localhost:56902
	worldProxyTcpIp.Timeout = 50000 * time.Millisecond

	gui.RunGui(&worldProxyTcpIp)
}
