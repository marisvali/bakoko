package main

import (
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
)

var w World

func main() {
	frameIdx := 0
	player1 := PlayerProxy{}
	player1.Endpoint = "localhost:56901"
	player2 := PlayerProxy{}
	player2.Endpoint = "localhost:56902"
	guiProxy := GuiProxy{}
	guiProxy.Endpoint = "localhost:56903"
	var watcher FolderWatcher
	watcher.Folder = "world-data"

	for w.Over.Eq(I(0)) {
		if watcher.FolderContentsChanged() {
			loadWorld(&w)
		}

		var input Input
		input.Player1Input = player1.GetInput() // Should block.
		input.Player2Input = player2.GetInput() // Should block.

		if input.Player1Input.Reload || input.Player2Input.Reload {
			loadWorld(&w)
		}
		w.Step(&input, frameIdx)

		guiProxy.SendPaintData(&w.DebugInfo) // Should not block.
		player1.SendWorld(&w)                // Should not block.
		player2.SendWorld(&w)                // Should not block.

		if input.Player1Input.Quit || input.Player2Input.Quit {
			break
		}
		frameIdx++
		w.JustReloaded = ZERO
	}
}

type worldData struct {
	BallSpeed                int
	BallDec                  int
	BallDiameter             int
	Player1X                 int
	Player1Y                 int
	Player1Speed             int
	Player1Health            int
	Player1NBalls            int
	Player1BallType          int
	Player1Diameter          int
	Player1StunnedImobilizes bool
	Player2X                 int
	Player2Y                 int
	Player2Speed             int
	Player2Health            int
	Player2NBalls            int
	Player2BallType          int
	Player2Diameter          int
	Player2StunnedImobilizes bool
	ObstacleSize             int
	Level                    string
}

func loadWorld(w *World) {
	data := loadWorldData("world-data")

	w.BallSpeed = I(data.BallSpeed)
	w.BallDec = I(data.BallDec)
	w.BallDiameter = I(data.BallDiameter)
	w.Player1.Bounds.Center.X = I(data.Player1X)
	w.Player1.Bounds.Center.Y = I(data.Player1Y)
	w.Player1.Speed = I(data.Player1Speed)
	w.Player1.Health = I(data.Player1Health)
	w.Player1.NBalls = I(data.Player1NBalls)
	w.Player1.BallType = I(data.Player1BallType)
	w.Player1.Bounds.Diameter = I(data.Player1Diameter)
	w.Player1.StunnedImobilizes = data.Player1StunnedImobilizes
	w.Player2.Bounds.Center.X = I(data.Player2X)
	w.Player2.Bounds.Center.Y = I(data.Player2Y)
	w.Player2.Speed = I(data.Player2Speed)
	w.Player2.Health = I(data.Player2Health)
	w.Player2.NBalls = I(data.Player2NBalls)
	w.Player2.BallType = I(data.Player2BallType)
	w.Player2.Bounds.Diameter = I(data.Player2Diameter)
	w.Player2.StunnedImobilizes = data.Player2StunnedImobilizes
	w.ObstacleSize = I(data.ObstacleSize)
	levelString := ReadAllText(data.Level)
	var balls1 []Pt
	//var balls2 []Pt
	w.Obstacles, balls1, _ = LevelFromString(levelString)
	w.Balls = []Ball{} // reset balls
	for i := range balls1 {
		b := Ball{
			//Pos:            Pt{player.Pos.X + (player.Diameter+30*Unit)/2 + 2*Unit, player.Pos.Y},
			Bounds: Circle{
				Center: Pt{balls1[i].X.Times(w.ObstacleSize).Plus(w.ObstacleSize.DivBy(TWO)),
					balls1[i].Y.Times(w.ObstacleSize).Plus(w.ObstacleSize.DivBy(TWO))},
				Diameter: w.BallDiameter},
			MoveDir:        IPt(0, 0),
			Speed:          ZERO,
			CanBeCollected: false,
			Type:           w.Player1.BallType,
		}
		w.Balls = append(w.Balls, b)
	}
	w.JustReloaded = ONE
}

func loadWorldData(folder string) (data worldData) {
	// Read from the disk over and over until a full read is possible.
	// This repetition is meant to avoid crashes due to reading files
	// while they are still being written.
	// It's a hack but possibly a quick and very useful one.
	CheckCrashes = false
	for {
		CheckFailed = nil
		LoadJSON(folder+"/world.json", &data)
		if CheckFailed == nil {
			break
		}
	}
	CheckCrashes = true
	return
}
