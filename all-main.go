package main

import (
	. "playful-patterns.com/bakoko/ai"
	. "playful-patterns.com/bakoko/gui"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
	. "playful-patterns.com/bakoko/world/world-run"
)

func main() {
	var w World
	recordingFile := "recorded-inputs-01"

	player1 := PlayerProxyRegular{}
	player2 := PlayerProxyRegular{}
	guiProxy := GuiProxyRegular{} // This isn't used yet.

	worldProxy1 := WorldProxyRegular{}
	worldProxy2 := WorldProxyRegular{}

	playerInputChannel1 := make(chan []byte)
	worldChannel1 := make(chan []byte)
	playerInputChannel2 := make(chan []byte)
	worldChannel2 := make(chan []byte)

	player1.InputChannel = playerInputChannel1
	worldProxy1.InputChannel = playerInputChannel1
	player1.WorldChannel = worldChannel1
	worldProxy1.WorldChannel = worldChannel1

	player2.InputChannel = playerInputChannel2
	worldProxy2.InputChannel = playerInputChannel2
	player2.WorldChannel = worldChannel2
	worldProxy2.WorldChannel = worldChannel2

	go RunWorld(&w, &player1, &player2, &guiProxy, recordingFile)
	go RunAi(&guiProxy, &worldProxy2)
	RunGui(&worldProxy1)
}
