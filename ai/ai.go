package main

import (
	"image/color"
	"os"
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
	"time"
)

type PlayerAI struct {
	TargetPt          Pt
	HasTarget         bool
	GuiProxy          GuiProxy
	DebugInfo         DebugInfo
	PauseBetweenShots time.Duration
	LastShot          time.Time
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

	if w.JustReloaded.Eq(ONE) {
		// Reset target if the world just reloaded.
		p.HasTarget = false
	}

	// TODO: find a more generic way of selecting which player is which.
	player := &w.Player2

	if player.Stunned.Gt(ZERO) {
		return
	}

	if time.Now().Sub(p.LastShot) > p.PauseBetweenShots {
		ballStart := player.Bounds.Center
		ballEnd := w.Player1.Bounds.Center
		if pathIsClear(w.Obstacles, w.ObstacleSize, ballStart, ballEnd, U(50)) {
			input.Shoot = true
			input.ShootPt = ballEnd
			p.LastShot = time.Now()
		}
	}

	//return

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

		p.DebugInfo.Points = []DebugPoint{}
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
				p.DebugInfo.Points = append(p.DebugInfo.Points, pt)
			}
		}

		var pathfinding Pathfinding
		pathfinding.Initialize(mw)
		path := pathfinding.FindPath(startPt, endPt)
		// Transform path coordinates into world-main coordinates.
		var pathWorld []Pt
		for _, pt := range path {
			pathWorld = append(pathWorld, pt.Times(sizeW))
			p.DebugInfo.Points = append(p.DebugInfo.Points, DebugPoint{pathWorld[len(pathWorld)-1], U(10), color.RGBA{255, 123, 0, 255}})
		}

		p.DebugInfo.Points = append(p.DebugInfo.Points, DebugPoint{Pt{startPt.X.Times(sizeW), startPt.Y.Times(sizeW)}, U(10), color.RGBA{255, 255, 255, 255}})
		p.DebugInfo.Points = append(p.DebugInfo.Points, DebugPoint{Pt{endPt.X.Times(sizeW), endPt.Y.Times(sizeW)}, U(10), color.RGBA{0, 0, 0, 255}})

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
		oldInput := input
		input = MoveStraightLine(player.Bounds.Center, p.TargetPt)
		input.Shoot = oldInput.Shoot
		input.ShootPt = oldInput.ShootPt
	}

	return
}

func pathIsClear(obstacles Matrix, obstacleSize Int, start Pt, end Pt, ballSize Int) bool {
	squares := obstaclesToSquares(obstacles, obstacleSize)
	// Check if we can travel to newPos without collision.
	// CircleSquareCollision doesn't return oldPos as a collision point.
	intersects, _, _ :=
		CircleSquaresCollision(start, end, ballSize, squares)
	return !intersects
}

func main() {
	var worldProxy WorldProxy
	var guiProxy GuiProxy
	var w World
	var ai PlayerAI
	ai.PauseBetweenShots = 1500 * time.Millisecond
	ai.LastShot = time.Now()

	worldProxy.Endpoint = os.Args[1] // localhost:56901 or localhost:56902
	worldProxy.Timeout = 0 * time.Millisecond
	guiProxy.Endpoint = os.Args[2]

	for {
		input := ai.Step(&w)

		// This should block as the AI doesn't make sense if it doesn't
		// synchronize with the simulation.
		for {
			if err := worldProxy.Connect(); err != nil {
				continue // Retry from the beginning.
			}

			if err := worldProxy.SendInput(&input); err != nil {
				continue // Retry from the beginning.
			}

			if err := worldProxy.GetWorld(&w); err != nil {
				continue // Retry from the beginning.
			}

			break
		}

		// This may or may not block, who cares?
		//guiProxy.SendPaintData(&ai.DebugInfo)
	}
}

func obstaclesToSquares(obstacles Matrix, obstacleSize Int) (squares []Square) {
	for row := I(0); row.Lt(obstacles.NRows()); row.Inc() {
		for col := I(0); col.Lt(obstacles.NCols()); col.Inc() {
			if obstacles.Get(row, col).Neq(I(0)) {
				half := obstacleSize.DivBy(I(2))
				squares = append(squares, Square{
					Center: Pt{col.Times(obstacleSize).Plus(half), row.Times(obstacleSize).Plus(half)},
					Size:   obstacleSize,
				})
			}
		}
	}
	return
}

func GetMatrixPointClosestToWorld(m Matrix, size Int, offset Pt, pos Pt) Pt {
	pos2 := pos.Minus(offset)
	// Get point by doing integer division (rounding down).
	first := pos2.DivBy(size)

	// The point in the world is surrounded by 4 points in the matrix.
	options := []Pt{
		first,
		{first.X.Plus(ONE), first.Y},
		{first.X, first.Y.Plus(ONE)},
		{first.X.Plus(ONE), first.Y.Plus(ONE)},
	}

	// Transform these options back into world points.
	optionsWorld := make([]Pt, 4)
	for i, pt := range options {
		worldPt := pt.Times(size).Plus(offset)
		optionsWorld[i] = worldPt
	}

	// Find the distances between each of the 4 points and the target point.
	distances := make([]Int, 4)
	for i, pt := range optionsWorld {
		distances[i] = pt.SquaredDistTo(pos)
	}

	// I don't want to spend time learning how to do a proper sorting here.
	// So, do it the brute-force way.
	for {
		// Get the smallest distance that's still valid.
		minIdx := -1
		for i := range distances {
			if distances[i].Gt(ZERO) && (minIdx < 0 || distances[i].Lt(distances[minIdx])) {
				minIdx = i
			}
		}
		// No more distances available. It means all points are disabled
		// on the matrix. Just return the first one.
		if minIdx == -1 {
			return options[0]
		}

		if m.Get(options[minIdx].Y, options[minIdx].X).Eq(ZERO) {
			return options[minIdx]
		} else {
			distances[minIdx] = I(-1)
		}
	}
}
