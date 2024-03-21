package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"os"
	. "playful-patterns.com/bakoko/gui"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
	"time"
)

func main() {
	var worldProxyTcpIp WorldProxyTcpIp
	worldProxyTcpIp.Endpoint = os.Args[1] // localhost:56901 or localhost:56902
	worldProxyTcpIp.Timeout = 50000 * time.Millisecond

	RunGuiSplitPlay(&worldProxyTcpIp)
}

func RunGuiSplitPlay(worldProxy WorldProxy) {
	var g Gui
	g.Init(worldProxy, nil, nil, "")

	// Start the game.
	err := ebiten.RunGame(&g)
	Check(err)
}
