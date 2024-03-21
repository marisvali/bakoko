package main

import (
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
	. "playful-patterns.com/bakoko/world/world-run"
)

func main() {
	mainRecord()
}

func mainRecord() {
	player1 := PlayerProxyTcpIp{}
	player1.Endpoint = "localhost:56901"
	player2 := PlayerProxyTcpIp{}
	player2.Endpoint = "localhost:56902"
	guiProxy := GuiProxyTcpIp{}
	guiProxy.Endpoint = "localhost:56903"

	RunWorldSplitPlay(&player1, &player2, &guiProxy, GetNewRecordingFile())
}

func RunWorldSplitPlay(player1 PlayerProxy, player2 PlayerProxy, guiProxy GuiProxy, recordingFile string) {
	var worldRunner WorldRunner
	worldRunner.Initialize(player1, player2, recordingFile)
	for {
		// First, send the current world to players and get their reactions.
		var input Input
		input.Player1Input = player1.SendWorldGetInput(&worldRunner.w) // Blocks.
		input.Player2Input = player2.SendWorldGetInput(&worldRunner.w) // Blocks.

		// Second, use their reactions to update the world.
		worldRunner.Step(input)
	}
}
