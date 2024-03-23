package world_run

import (
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/world"
)

type WorldRunner struct {
	w             World
	frameIdx      int
	watcher       FolderWatcher
	recordingFile string
	currentInputs []PlayerInput
}

func (wr *WorldRunner) Initialize(recordingFile string, folderWatchingEnabled bool) {
	wr.frameIdx = 0
	if folderWatchingEnabled {
		wr.watcher.Folder = "world-data"
	}
	wr.recordingFile = recordingFile
	LoadWorld(&wr.w)
}

func (wr *WorldRunner) Step(input Input) {
	if wr.recordingFile != "" {
		wr.currentInputs = append(wr.currentInputs, input.Player1Input)
		SerializeInputs(wr.currentInputs, wr.recordingFile)
	}

	//if input.Player1Input.Quit || input.Player2Input.Quit {
	//	break
	//}

	wr.w.JustReloaded = ZERO
	if input.Player1Input.Reload || input.Player2Input.Reload || wr.watcher.FolderContentsChanged() {
		LoadWorld(&wr.w)
	}

	if !input.Player1Input.Pause && !input.Player2Input.Pause {
		wr.w.Step(&input, wr.frameIdx)
	}

	wr.frameIdx++
}

func (wr *WorldRunner) GetWorld() *World {
	return &wr.w
}

func (wr *WorldRunner) GetDebugInfo() *DebugInfo {
	return &wr.w.DebugInfo
}
