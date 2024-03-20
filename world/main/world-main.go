package main

import (
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world/world-run"
)

func main() {
	mainRecord()
}

func mainRecord() {
	recordingFile := "recorded-inputs-01"

	player1 := PlayerProxyTcpIp{}
	player1.Endpoint = "localhost:56901"
	player2 := PlayerProxyTcpIp{}
	player2.Endpoint = "localhost:56902"
	guiProxy := GuiProxyTcpIp{}
	guiProxy.Endpoint = "localhost:56903"

	RunWorldRecord(&player1, &player2, &guiProxy, recordingFile)
}

//func mainReplay() {
//	playbackFile := "recorded-inputs-01"
//	playbackInputs := DeserializeInputs(playbackFile)
//
//	frameIdx := 0
//	player1 := PlayerProxy{}
//	player1.Endpoint = "localhost:56901"
//	player2 := PlayerProxy{}
//	player2.Endpoint = "localhost:56902"
//	guiProxy := GuiProxy{}
//	guiProxy.Endpoint = "localhost:56903"
//	var watcher FolderWatcher
//	watcher.Folder = "world-data"
//
//	for w.Over.Eq(I(0)) {
//		if watcher.FolderContentsChanged() {
//			loadWorld(&w)
//		}
//
//		var input Input
//		input.Player1Input = player1.GetInput()
//		if frameIdx < len(playbackInputs) {
//			input.Player1Input = playbackInputs[frameIdx]
//		}
//		input.Player2Input = player2.GetInput() // Should block.
//
//		if input.Player1Input.Reload || input.Player2Input.Reload {
//			loadWorld(&w)
//		}
//
//		if !input.Player1Input.Pause && !input.Player2Input.Pause {
//			w.Step(&input, frameIdx)
//		}
//
//		guiProxy.SendPaintData(&w.DebugInfo) // Should not block.
//		player1.SendWorld(&w)                // Should not block.
//		player2.SendWorld(&w)                // Should not block.
//
//		if input.Player1Input.Quit || input.Player2Input.Quit {
//			break
//		}
//		frameIdx++
//		w.JustReloaded = ZERO
//	}
//}
