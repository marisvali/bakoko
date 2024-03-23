package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"os"
	. "playful-patterns.com/bakoko/ai"
	. "playful-patterns.com/bakoko/gui"
	. "playful-patterns.com/bakoko/world"
	. "playful-patterns.com/bakoko/world/world-run"
)

func main() {
	if len(os.Args) == 1 {
		RunGuiFusedPlay(GetNewRecordingFile())
	} else {
		RunGuiFusedPlayback(os.Args[1])
		//RunGuiFusedPlayback("d:/gms/bakoko/recordings/recorded-inputs-2024-03-20-000000")
	}
}

func RunGuiFusedPlay(recordingFile string) {
	var worldRunner WorldRunner
	var player2Ai PlayerAI
	worldRunner.Initialize(recordingFile, true)

	var g Gui
	g.Init(nil, &worldRunner, &player2Ai, "", []string{})

	// Start the game.
	err := ebiten.RunGame(&g)
	Check(err)
}

func RunGuiFusedPlayback(recordingFile string) {
	var worldRunner WorldRunner
	var player2Ai PlayerAI
	worldRunner.Initialize("", false)

	var g Gui
	g.Init(nil, &worldRunner, &player2Ai, recordingFile, []string{})

	// Start the game.
	err := ebiten.RunGame(&g)
	Check(err)
}
