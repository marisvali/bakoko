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
	worldProxy1.PlayerProxy = &player1
	worldProxy2 := WorldProxyRegular{}
	worldProxy2.PlayerProxy = &player2

	player1.WorldProxy = &worldProxy1
	player2.WorldProxy = &worldProxy2
	//
	//playerInputChannel1 := make(chan []byte)
	//playerInputChannel2 := make(chan []byte)

	go RunWorld(&w, &player1, &player2, &guiProxy, recordingFile)
	go RunAi(&guiProxy, &worldProxy2)
	RunGui(&worldProxy1)
}
