package main

import . "playful-patterns.com/bakoko"

func main() {
	var input Input
	var w World
	//if utils.FileExists("world.bin") {
	//	w.DeserializeFromFile("world.bin")
	//}
	w.Player1 = Character{
		Pos:      Point{X: 60 * Unit, Y: 260 * Unit},
		Diameter: 50 * Unit,
		NBalls:   3}
	w.Player2 = Character{
		Pos:      Point{X: 180 * Unit, Y: 60 * Unit},
		Diameter: 50 * Unit,
		NBalls:   3}
	w.Balls = []Ball{
		{
			Pos:      Point{X: 120 * Unit, Y: 70 * Unit},
			Diameter: 30 * Unit,
		},
		{
			Pos:      Point{X: 90 * Unit, Y: 210 * Unit},
			Diameter: 30 * Unit,
		},
		{
			Pos:      Point{X: 190 * Unit, Y: 140 * Unit},
			Diameter: 30 * Unit,
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
