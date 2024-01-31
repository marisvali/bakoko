package main

import (
	"image/color"
	"os"
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
)

type PlayerAI struct {
	TargetPt  Pt
	HasTarget bool
}

func PlayerIsAt(p *Player, pt Pt) bool {
	return p.Bounds.Center.SquaredDistTo(pt).Lt(U(5).Sqr())
}
func (p *PlayerAI) Step(w *World) (input PlayerInput) {
	// Check somehow if the world-main is initialized.
	if w.Obstacles.NRows().Leq(ZERO) {
		// If there's no world-main matrix, we can probably safely assume
		// the world-main is not initialized.
		return
	}

	// TODO: find a more generic way of selecting which player is which.
	player := &w.Player2
	finalTarget := w.Player1.Bounds.Center

	if p.HasTarget {
		// If we're at the target, disable the target which signals we need
		// a new path.
		if PlayerIsAt(player, p.TargetPt) {
			p.HasTarget = false
		}
	}

	// Only compute a new path when we don't have a target to go to.
	if !p.HasTarget && !PlayerIsAt(player, finalTarget) {
		mw, sizeW, offsetW := GetWalkableMatrix(w.Obstacles, w.ObstacleSize, player.Bounds.Diameter)

		playerPos := player.Bounds.Center
		startPt := GetMatrixPointClosestToWorld(mw, sizeW, offsetW, playerPos)
		endPt := GetMatrixPointClosestToWorld(mw, sizeW, offsetW, finalTarget)

		w.DebugPoints = []DebugPoint{}
		//w.Obs = []Square{
		//	{Pt{U(100), U(100)}, U(10)},
		//	{playerPos, U(10)},
		//	{finalTarget, U(10)},
		//}
		for y := I(0); y.Lt(mw.NRows()); y.Inc() {
			for x := I(0); x.Lt(mw.NCols()); x.Inc() {
				var pt DebugPoint
				pt.Pos = Pt{x.Times(sizeW), y.Times(sizeW)}
				pt.Pos.Add(offsetW)
				pt.Size = U(3)
				if mw.Get(y, x).Eq(ZERO) {
					pt.Col = color.RGBA{0, 0, 255, 255}
				} else {
					pt.Col = color.RGBA{255, 0, 0, 255}
				}
				w.DebugPoints = append(w.DebugPoints, pt)
			}
		}

		var pathfinding Pathfinding
		pathfinding.Initialize(mw)
		path := pathfinding.FindPath(startPt, endPt)
		// Transform path coordinates into world-main coordinates.
		var pathWorld []Pt
		for _, pt := range path {
			pathWorld = append(pathWorld, pt.Times(sizeW))
			w.DebugPoints = append(w.DebugPoints, DebugPoint{pathWorld[len(pathWorld)-1], U(10), color.RGBA{255, 123, 0, 255}})
		}

		w.DebugPoints = append(w.DebugPoints, DebugPoint{Pt{startPt.X.Times(sizeW), startPt.Y.Times(sizeW)}, U(10), color.RGBA{255, 255, 255, 255}})
		w.DebugPoints = append(w.DebugPoints, DebugPoint{Pt{endPt.X.Times(sizeW), endPt.Y.Times(sizeW)}, U(10), color.RGBA{255, 123, 255, 255}})

		if len(pathWorld) == 0 {
			// Don't do anything.
			p.HasTarget = false
		} else {
			// Remove the first point if the player is basically there.
			if pathWorld[0].SquaredDistTo(playerPos).Lt(U(5).Sqr()) {
				pathWorld = pathWorld[1:]
			}

			if len(pathWorld) == 0 {
				// Just move towards the target.
				p.HasTarget = true
				p.TargetPt = finalTarget
			} else {
				// Move towards the next point in the path
				p.HasTarget = true
				p.TargetPt = pathWorld[0]
			}
		}
	}

	if p.HasTarget {
		input = MoveStraightLine(player.Bounds.Center, p.TargetPt)
	}

	return
}

func main() {
	var peer SimulationPeer
	var w World
	var ai PlayerAI

	peer.Endpoint = os.Args[1] // localhost:56901 or localhost:56902

	for {
		input := ai.Step(&w)
		peer.SendInput(&input)
		peer.GetWorld(&w)
	}
}
