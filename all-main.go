package main

import (
	"fmt"
	"os"
	. "playful-patterns.com/bakoko/gui"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
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
	worldProxy := WorldProxyPlayback{}
	RunGuiPlayback(&worldProxy, recordingFile)
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
	recordingFile := getNewRecordingFile()
	worldProxy := WorldProxyPlayback{}
	worldProxy.RecordingFile = recordingFile
	RunGui(&worldProxy)
}
