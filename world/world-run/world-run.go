package world_run

import (
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
)

type WorldRunner struct {
	w             World
	frameIdx      int
	watcher       FolderWatcher
	recordingFile string
	currentInputs []PlayerInput
	player1       PlayerProxy
	player2       PlayerProxy
}

func (wr *WorldRunner) Initialize(player1 PlayerProxy, player2 PlayerProxy, recordingFile string) {
	wr.frameIdx = 0
	wr.watcher.Folder = "world-data"
	wr.recordingFile = recordingFile
	wr.player1 = player1
	wr.player2 = player2
	wr.player1.SendWorld(&wr.w) // Should not block.
	wr.player2.SendWorld(&wr.w) // Should not block.
}

func (wr *WorldRunner) Step() {
	var input Input
	input.Player1Input = wr.player1.GetInput() // Should block.
	input.Player2Input = wr.player2.GetInput() // Should block.

	if wr.recordingFile != "" {
		wr.currentInputs = append(wr.currentInputs, input.Player1Input)
		SerializeInputs(wr.currentInputs, wr.recordingFile)
	}

	// Only change the world in this well defined part. This way, the players
	// can get pointers to our world data and use them without making copies
	// of the world, because the world remains unchanged between when we send
	// the world to a player and we get the input from that player.
	// !!! Start changing the world.
	wr.w.JustReloaded = ZERO
	if input.Player1Input.Reload || input.Player2Input.Reload || wr.watcher.FolderContentsChanged() {
		LoadWorld(&wr.w)
	}

	if !input.Player1Input.Pause && !input.Player2Input.Pause {
		wr.w.Step(&input, wr.frameIdx)
	}
	// !!! Stop changing the world.

	//guiProxy.SendPaintData(&world.DebugInfo) // Should not block.
	wr.player1.SendWorld(&wr.w) // Should not block.
	wr.player2.SendWorld(&wr.w) // Should not block.

	//if input.Player1Input.Quit || input.Player2Input.Quit {
	//	break
	//}
	wr.frameIdx++
}

func RunWorldRecord(player1 PlayerProxy, player2 PlayerProxy, guiProxy GuiProxy, recordingFile string) {
	var worldRunner WorldRunner
	worldRunner.Initialize(player1, player2, recordingFile)
	for worldRunner.w.Over.Eq(I(0)) {
		worldRunner.Step()
	}
}
