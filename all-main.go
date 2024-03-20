package main

import (
	"fmt"
	"os"
	. "playful-patterns.com/bakoko/gui"
	. "playful-patterns.com/bakoko/world"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		RunGuiFusedPlay(getNewRecordingFile())
	} else {
		//RunGuiFusedPlayback(os.Args[1])
		RunGuiFusedPlayback("d:/gms/bakoko/recordings/recorded-inputs-2024-03-20-000000")
	}
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
