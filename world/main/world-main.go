package main

import (
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
	. "playful-patterns.com/bakoko/world/world-run"
)

// 3 possible run modes: FusedRecording, FusedPlayback, SplitRecording
// Run the world in SplitRecording mode.
func main() {
	player1 := PlayerProxyTcpIp{}
	player1.Endpoint = "localhost:56901"
	player2 := PlayerProxyTcpIp{}
	player2.Endpoint = "localhost:56902"
	guiProxy := GuiProxyTcpIp{}
	guiProxy.Endpoint = "localhost:56903"

	var worldRunner WorldRunner
	worldRunner.Initialize(GetNewRecordingFile())
	for {
		// First, send the current world to players and get their reactions.
		var input Input
		input.Player1Input = *player1.SendWorldGetInput(worldRunner.GetWorld()) // Blocks.
		input.Player2Input = *player2.SendWorldGetInput(worldRunner.GetWorld()) // Blocks.

		// Second, use their reactions to update the world.
		worldRunner.Step(input)
	}
}
