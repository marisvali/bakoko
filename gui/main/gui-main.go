package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"os"
	. "playful-patterns.com/bakoko/gui"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
	"time"
)

// 3 possible run modes: FusedRecording, FusedPlayback, SplitRecording
// Run the world in SplitRecording mode.
func main() {
	var worldProxyTcpIp WorldProxyTcpIp
	worldProxyTcpIp.Endpoint = os.Args[1] // localhost:56901 or localhost:56902
	worldProxyTcpIp.Timeout = 50000 * time.Millisecond

	var g Gui
	g.Init(&worldProxyTcpIp, nil, nil, "")

	// Start the game.
	err := ebiten.RunGame(&g)
	Check(err)
}
