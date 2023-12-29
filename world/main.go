package main

import (
	. "playful-patterns.com/bakoko"
)
import . "playful-patterns.com/bakoko/ints"

func main() {
	var input Input
	var w World
	//if utils.FileExists("world.bin") {
	//	w.DeserializeFromFile("world.bin")
	//}
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

	frameIdx := 0
	for {
		//frameStart := time.Now()
		WaitForFile("input-ready")
		input.DeserializeFromFile("input.bin")
		w.Step(&input, frameIdx)
		w.SerializeToFile("world.bin")
		DeleteFile("input-ready")
		TouchFile("world-ready")

		//frameDuration := time.Since(frameStart)
		//fmt.Printf("frame idx: %d duration: %d\n", frameIdx, frameDuration.Milliseconds())
		frameIdx++
	}
}
