package main

import (
	"fmt"
	"math/rand"
	. "playful-patterns.com/bakoko"
	"time"
)
import . "playful-patterns.com/bakoko/ints"

var input Input
var w World

func init() {
	w.Player1 = Player{
		Pos:      UPt(60, 260),
		Diameter: U(50),
		NBalls:   I(3)}
	w.Player2 = Player{
		Pos:      UPt(180, 60),
		Diameter: U(50),
		NBalls:   I(3)}
	w.Balls = []Ball{
		{
			Pos:      UPt(120, 70),
			Diameter: U(30),
		},
		{
			Pos:      UPt(90, 210),
			Diameter: U(30),
		},
		{
			Pos:      UPt(190, 140),
			Diameter: U(30),
		}}
}
func main5() {
	originalWorld := w
	w.SerializeToFile("world.bin")
	inputs := DeserializeInputs("recorded-inputs")
	start := time.Now()
	bigIdx := 0
	for ; bigIdx < 1000000; bigIdx++ {
		frameIdx := 0
		w := originalWorld
		for ; !input.Quit; frameIdx++ {
			input = Input{}
			if frameIdx < len(inputs) {
				input = inputs[frameIdx]
			}

			w.Step(&input, frameIdx)
		}
	}
	fmt.Println(time.Since(start).Seconds())
	fmt.Printf("%.12f\n", time.Since(start).Seconds()/float64(bigIdx))
	w.SerializeToFile("world.bin")
}
func main4() {
	w.SerializeToFile("world.bin")
	inputs := DeserializeInputs("recorded-inputs")

	frameIdx := 0
	for ; !input.Quit; frameIdx++ {
		//frameStart := time.Now()
		WaitForFile("input-ready")
		//input.DeserializeFromFile("input.bin")
		input = Input{}
		if frameIdx < len(inputs) {
			input = inputs[frameIdx]
		}

		w.Step(&input, frameIdx)
		w.SerializeToFile("world.bin")
		DeleteFile("input-ready")
		TouchFile("world-ready")

		//frameDuration := time.Since(frameStart)
		//fmt.Printf("frame idx: %d duration: %d\n", frameIdx, frameDuration.Milliseconds())
	}
}
func main3() {
	w.SerializeToFile("world.bin")

	frameIdx := 0
	for ; !input.Quit; frameIdx++ {
		//frameStart := time.Now()
		WaitForFile("input-ready")
		//input.DeserializeFromFile("input.bin")
		input = Input{}

		input.MoveLeft = rand.Int()%6 == 0
		input.MoveRight = rand.Int()%5 == 0
		input.MoveUp = rand.Int()%6 == 0
		input.MoveDown = rand.Int()%6 == 0
		input.Quit = frameIdx == 3600

		w.Step(&input, frameIdx)
		w.SerializeToFile("world.bin")
		DeleteFile("input-ready")
		TouchFile("world-ready")

		//frameDuration := time.Since(frameStart)
		//fmt.Printf("frame idx: %d duration: %d\n", frameIdx, frameDuration.Milliseconds())
	}
}

var currentInputs []Input

func GetNewRecordingFile() string {
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

func main() {
	//if utils.FileExists("world.bin") {
	//	w.DeserializeFromFile("world.bin")
	//}

	playbackFile := ""
	// Change the file name to choose a specific playthrough to play back or
	// comment this line disable playback.
	playbackFile = "recordings/recorded-inputs-2023-12-31-000000"

	var playbackInputs []Input
	if playbackFile != "" {
		playbackInputs = DeserializeInputs(playbackFile)
	}

	w.SerializeToFile("world.bin")
	recordingFile := GetNewRecordingFile()
	frameIdx := 0
	for ; !input.Quit && (playbackFile == "" || frameIdx < len(playbackInputs)); frameIdx++ {
		WaitForFile("input-ready")
		if playbackFile != "" {
			input = playbackInputs[frameIdx]
		} else {
			input.DeserializeFromFile("input.bin")
		}
		currentInputs = append(currentInputs, input)
		SerializeInputs(currentInputs, recordingFile)

		w.Step(&input, frameIdx)
		w.SerializeToFile("world.bin")
		DeleteFile("input-ready")
		TouchFile("world-ready")
	}
}
