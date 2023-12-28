package main

import (
	"playful-patterns.com/bakoko/utils"
	"playful-patterns.com/bakoko/world"
)

func main() {
	var input world.Input
	var w world.World
	//if utils.FileExists("world.bin") {
	//	w.DeserializeFromFile("world.bin")
	//}
	w.Player1 = world.Character{
		Pos:      world.Point{X: 60 * 1000, Y: 260 * 1000},
		Diameter: 50 * 1000,
		NBalls:   3}
	w.Player2 = world.Character{
		Pos:      world.Point{X: 180 * 1000, Y: 60 * 1000},
		Diameter: 50 * 1000,
		NBalls:   3}
	w.Balls = []world.Ball{
		{
			Pos:      world.Point{X: 120 * 1000, Y: 70 * 1000},
			Diameter: 30 * 1000,
		},
		{
			Pos:      world.Point{X: 90 * 1000, Y: 210 * 1000},
			Diameter: 30 * 1000,
		},
		{
			Pos:      world.Point{X: 190 * 1000, Y: 140 * 1000},
			Diameter: 30 * 1000,
		}}

	frameIdx := 0
	for {
		//frameStart := time.Now()
		utils.WaitForFile("input-ready")
		input.DeserializeFromFile("input.bin")
		w.Step(&input, frameIdx)
		w.SerializeToFile("world.bin")
		utils.DeleteFile("input-ready")
		utils.TouchFile("world-ready")

		//frameDuration := time.Since(frameStart)
		//fmt.Printf("frame idx: %d duration: %d\n", frameIdx, frameDuration.Milliseconds())
		frameIdx++
	}
}
