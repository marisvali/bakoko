package world_run

import (
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
)

func RunWorldRecord(w *World, player1 PlayerProxy, player2 PlayerProxy, guiProxy GuiProxy, recordingFile string) {
	frameIdx := 0
	var currentInputs []PlayerInput
	var watcher FolderWatcher
	watcher.Folder = "world-data"
	for w.Over.Eq(I(0)) {
		if watcher.FolderContentsChanged() {
			LoadWorld(w)
		}

		var input Input
		input.Player1Input = player1.GetInput() // Should block.
		input.Player2Input = player2.GetInput() // Should block.

		currentInputs = append(currentInputs, input.Player1Input)
		SerializeInputs(currentInputs, recordingFile)

		if input.Player1Input.Reload || input.Player2Input.Reload {
			LoadWorld(w)
		}

		if !input.Player1Input.Pause && !input.Player2Input.Pause {
			w.Step(&input, frameIdx)
		}

		guiProxy.SendPaintData(&w.DebugInfo) // Should not block.
		player1.SendWorld(w)                 // Should not block.
		player2.SendWorld(w)                 // Should not block.

		if input.Player1Input.Quit || input.Player2Input.Quit {
			break
		}
		frameIdx++
		w.JustReloaded = ZERO
	}
}

func RunWorldPlayback(w *World, player1 PlayerProxy, player2 PlayerProxy, guiProxy GuiProxy, playbackFile string) {
	playbackInputs := DeserializeInputs(playbackFile)

	frameIdx := 0
	var currentInputs []PlayerInput
	var watcher FolderWatcher
	watcher.Folder = "world-data"
	for w.Over.Eq(I(0)) {
		if watcher.FolderContentsChanged() {
			LoadWorld(w)
		}

		var input Input
		input.Player1Input = player1.GetInput() // Should block.
		if frameIdx < len(playbackInputs) {
			input.Player1Input = playbackInputs[frameIdx]
		}
		input.Player2Input = player2.GetInput() // Should block.

		currentInputs = append(currentInputs, input.Player1Input)

		if input.Player1Input.Reload || input.Player2Input.Reload {
			LoadWorld(w)
		}

		if !input.Player1Input.Pause && !input.Player2Input.Pause {
			w.Step(&input, frameIdx)
		}

		guiProxy.SendPaintData(&w.DebugInfo) // Should not block.
		player1.SendWorld(w)                 // Should not block.
		player2.SendWorld(w)                 // Should not block.

		if input.Player1Input.Quit || input.Player2Input.Quit {
			break
		}
		frameIdx++
		w.JustReloaded = ZERO
	}
}
