package main

import (
	"fmt"
	"os"
	. "playful-patterns.com/bakoko/ai"
	. "playful-patterns.com/bakoko/gui"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
	. "playful-patterns.com/bakoko/world/world-run"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		mainRecord()
	} else {
		mainPlayback(os.Args[1])
	}
}

func mainPlayback(recordingFile string) {
	var w World

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

	go RunWorldPlayback(&w, &player1, &player2, &guiProxy, recordingFile)
	go RunAi(&guiProxy, &worldProxy2)
	RunGuiPlayback(&worldProxy1, recordingFile)
}

func getNewRecordingFile() string {
	date := time.Now()
	for i := 0; i < 1000000; i++ {
		filename := fmt.Sprintf("recordings/recorded-inputs-%04d-%02d-%02d-%06d",
			date.Year(), date.Month(), date.Day(), i)
		if !FileExists(filename) {
			return filename
		}
	}
	panic("Cannot record, no available filename found.")
}

func mainRecord() {
	var w World
	recordingFile := getNewRecordingFile()

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

	go RunWorldPlayback(&w, &player1, &player2, &guiProxy, recordingFile)
	go RunAi(&guiProxy, &worldProxy2)
	RunGui(&worldProxy1)
}
