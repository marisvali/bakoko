package ai

import (
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/world"
)

type PlayerAI struct {
	TargetPt                  Pt
	HasTarget                 bool
	DebugInfo                 DebugInfo
	PauseBetweenShots         Int
	LastShot                  Int
	walkableMatrix            Matrix
	sizeW                     Int
	offsetW                   Pt
	initializedWalkableMatrix bool
	pathfinding               Pathfinding
	frameIdx                  Int
}

func PlayerIsAt(p *Player, pt Pt) bool {
	return p.Bounds.Center.SquaredDistTo(pt).Lt(U(5).Sqr())
}

func (mind *PlayerAI) Initialize() {
	mind.HasTarget = false
	mind.initializedWalkableMatrix = false
	mind.PauseBetweenShots = I(90)
	mind.frameIdx = ZERO
	mind.LastShot = mind.frameIdx
}

func (mind *PlayerAI) Step(w *World) (input PlayerInput) {
	defer mind.frameIdx.Inc()

	// Check somehow if the world-main is initialized.
	if w.Obstacles.NRows().Leq(ZERO) {
		// If there's no world-main matrix, we can probably safely assume
		// the world-main is not initialized.
		return
	}

	// Check if us or the opponent is defeated.
	if w.Player1.Health.Eq(ZERO) || w.Player2.Health.Eq(ZERO) {
		return
	}

	// TODO: find a more generic way of selecting which body is which.
	body := &w.Player2

	if w.JustReloaded.Eq(ONE) {
		// Re-initialize the mind to initial conditions if the world just
		// reloaded. E.g. reset the target.
		mind.Initialize()
	}

	if !mind.initializedWalkableMatrix {
		mind.walkableMatrix, mind.sizeW, mind.offsetW = GetWalkableMatrix(w.Obstacles, w.ObstacleSize, body.Bounds.Diameter)
		mind.initializedWalkableMatrix = true
		mind.pathfinding.Initialize(mind.walkableMatrix)
	}

	if mind.frameIdx.Minus(mind.LastShot).Gt(mind.PauseBetweenShots) {
		ballStart := body.Bounds.Center
		ballEnd := w.Player1.Bounds.Center
		if pathIsClear(w, ballStart, ballEnd, U(50)) {
			input.Shoot = true
			input.ShootPt = ballEnd
			mind.LastShot = mind.frameIdx
		}
	}

	//return

	finalTarget := w.Player1.Bounds.Center
	if mind.HasTarget {
		// If we're at the target, disable the target which signals we need
		// a new path.
		if PlayerIsAt(body, mind.TargetPt) {
			mind.HasTarget = false
		}
	}

	// Only compute a new path when we don't have a target to go to.
	if !mind.HasTarget && !PlayerIsAt(body, finalTarget) {

		playerPos := body.Bounds.Center
		startPt := GetMatrixPointClosestToWorld(mind.walkableMatrix, mind.sizeW, mind.offsetW, playerPos)
		endPt := GetMatrixPointClosestToWorld(mind.walkableMatrix, mind.sizeW, mind.offsetW, finalTarget)

		//mind.DebugInfo.Points = []DebugPoint{}
		//w.Obs = []Square{
		//	{Pt{U(100), U(100)}, U(10)},
		//	{playerPos, U(10)},
		//	{finalTarget, U(10)},
		//}
		// Don't build debug info when not needed, as I'm trying to get max
		// execution time for replaying the world.
		//for y := I(0); y.Lt(mind.walkableMatrix.NRows()); y.Inc() {
		//	for x := I(0); x.Lt(mind.walkableMatrix.NCols()); x.Inc() {
		//		var pt DebugPoint
		//		pt.Pos = Pt{x.Times(mind.sizeW), y.Times(mind.sizeW)}
		//		pt.Pos.Add(mind.offsetW)
		//		pt.Size = U(3)
		//		if mind.walkableMatrix.Get(y, x).Eq(ZERO) {
		//			pt.Col = color.RGBA{0, 0, 255, 255}
		//		} else {
		//			pt.Col = color.RGBA{255, 0, 0, 255}
		//		}
		//		mind.DebugInfo.Points = append(mind.DebugInfo.Points, pt)
		//	}
		//}

		path := mind.pathfinding.FindPath(startPt, endPt)
		// Transform path coordinates into world-main coordinates.
		var pathWorld []Pt
		for _, pt := range path {
			pathWorld = append(pathWorld, pt.Times(mind.sizeW))
			//mind.DebugInfo.Points = append(mind.DebugInfo.Points, DebugPoint{pathWorld[len(pathWorld)-1], U(10), color.RGBA{255, 123, 0, 255}})
		}

		//mind.DebugInfo.Points = append(mind.DebugInfo.Points, DebugPoint{Pt{startPt.X.Times(mind.sizeW), startPt.Y.Times(mind.sizeW)}, U(10), color.RGBA{255, 255, 255, 255}})
		//mind.DebugInfo.Points = append(mind.DebugInfo.Points, DebugPoint{Pt{endPt.X.Times(mind.sizeW), endPt.Y.Times(mind.sizeW)}, U(10), color.RGBA{0, 0, 0, 255}})

		if len(pathWorld) == 0 {
			// Don't do anything.
			mind.HasTarget = false
		} else {
			// Remove the first point if the body is basically there.
			if pathWorld[0].SquaredDistTo(playerPos).Lt(U(5).Sqr()) {
				pathWorld = pathWorld[1:]
			}

			if len(pathWorld) == 0 {
				// Just move towards the target.
				mind.HasTarget = true
				mind.TargetPt = finalTarget
			} else {
				// Move towards the next point in the path
				mind.HasTarget = true
				mind.TargetPt = pathWorld[0]
			}
		}
	}

	if mind.HasTarget {
		moveInput := MoveStraightLine(body.Bounds.Center, mind.TargetPt)
		input.MoveLeft = moveInput.MoveLeft
		input.MoveRight = moveInput.MoveRight
		input.MoveUp = moveInput.MoveUp
		input.MoveDown = moveInput.MoveDown
	}

	return
}

func pathIsClear(w *World, start Pt, end Pt, ballSize Int) bool {
	squares := w.GetRelevantSquares(ballSize, start, end)
	// Check if we can travel to newPos without collision.
	// CircleSquareCollision doesn't return oldPos as a collision point.
	intersects, _, _ :=
		CircleSquaresCollision(start, end, ballSize, squares)
	return !intersects
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
