package world

import (
	"bytes"
	"fmt"
	"io"
	. "playful-patterns.com/bakoko/ints"
)

type Ball struct {
	Type           Int
	Bounds         Circle
	MoveDir        Pt
	Speed          Int
	CanBeCollected bool
}

type Player struct {
	Bounds   Circle
	NBalls   Int
	BallType Int
	Health   Int
	Speed    Int
}

type Matrix struct {
	cells []Int
	nRows Int
	nCols Int
}

func (m *Matrix) Clone() (c Matrix) {
	c.nRows = m.nRows
	c.nCols = m.nCols
	c.cells = append(c.cells, m.cells...)
	return
}

func (m *Matrix) Serialize(w io.Writer) {
	Serialize(w, m.nRows)
	Serialize(w, m.nCols)
	Serialize(w, m.cells)
}

func (m *Matrix) Deserialize(buf *bytes.Buffer) {
	Deserialize(buf, &m.nRows)
	Deserialize(buf, &m.nCols)
	m.cells = make([]Int, m.nRows.Times(m.nCols).ToInt64())
	Deserialize(buf, m.cells)
}

func (m *Matrix) Init(nRows, nCols Int) {
	m.nRows = nRows
	m.nCols = nCols
	m.cells = make([]Int, nRows.Times(nCols).ToInt64())
}

func (m *Matrix) Set(row, col, val Int) {
	m.cells[row.Times(m.nCols).Plus(col).ToInt64()] = val
}

func (m *Matrix) Get(row, col Int) Int {
	return m.cells[row.Times(m.nCols).Plus(col).ToInt64()]
}

func (m *Matrix) InBounds(pt Pt) bool {
	return pt.X.IsNonNegative() &&
		pt.Y.IsNonNegative() &&
		pt.Y.Lt(m.nRows) &&
		pt.X.Lt(m.nCols)
}

func (m *Matrix) NRows() Int {
	return m.nRows
}

func (m *Matrix) NCols() Int {
	return m.nCols
}

func (m *Matrix) PtToIndex(p Pt) Int {
	return p.Y.Times(m.nCols).Plus(p.X)
}

func (m *Matrix) IndexToPt(i Int) (p Pt) {
	p.X = i.Mod(m.nCols)
	p.Y = i.DivBy(m.nCols)
	return
}

type World struct {
	Player1      Player
	Player2      Player
	Balls        []Ball
	Over         Int
	Obstacles    Matrix
	ObstacleSize Int
	BallSpeed    Int
	BallDec      Int
	DebugInfo    DebugInfo
}

type PlayerInput struct {
	MoveLeft  bool
	MoveRight bool
	MoveUp    bool
	MoveDown  bool
	Shoot     bool
	ShootPt   Pt
	Quit      bool
	Reload    bool
}

type Input struct {
	Player1Input PlayerInput
	Player2Input PlayerInput
}

func (w *World) Serialize() []byte {
	buf := new(bytes.Buffer)
	Serialize(buf, w.Player1)
	Serialize(buf, w.Player2)
	SerializeSlice(buf, w.Balls)
	w.Obstacles.Serialize(buf)
	Serialize(buf, w.ObstacleSize)
	return buf.Bytes()
}

func (w *World) Deserialize(buf *bytes.Buffer) {
	Deserialize(buf, &w.Player1)
	Deserialize(buf, &w.Player2)
	DeserializeSlice(buf, &w.Balls)
	w.Obstacles.Deserialize(buf)
	Deserialize(buf, &w.ObstacleSize)
}

func ShootBall(player *Player, balls *[]Ball, pt Pt, ballSpeed Int) {
	if player.NBalls.Leq(I(0)) {
		return
	}

	moveDir := player.Bounds.Center.To(pt)
	//moveDir.SetLen(MU(6000))
	moveDir.SetLen(U(1))
	speed := ballSpeed

	ball := Ball{
		//Pos:            Pt{player.Pos.X + (player.Diameter+30*Unit)/2 + 2*Unit, player.Pos.Y},
		Bounds: Circle{
			Center:   player.Bounds.Center,
			Diameter: U(30)},
		MoveDir:        moveDir,
		Speed:          speed,
		CanBeCollected: false,
		Type:           player.BallType,
	}
	*balls = append(*balls, ball)
	// Infinite balls, for debugging purposes.
	//player.NBalls.Dec()
}

func ShootBallDebug(balls *[]Ball, orig, dest Pt, speed Int) {
	moveDir := orig.To(dest)
	moveDir.SetLen(U(1))

	ball := Ball{
		Bounds: Circle{
			Center:   orig,
			Diameter: U(30)},
		MoveDir:        moveDir,
		Speed:          speed,
		CanBeCollected: false,
		Type:           I(1),
	}
	*balls = append(*balls, ball)
}

// Compute the new position of a circle travelling from its current position
// along travelVec to a new position. Collisions are handled so that the final
// position takes into account any obstacles hit along the way.
// The logic of this function is that the circle travels for a length of
// travelLen in total and has no concept of time. So you can say it treats
// the movement as uniform, as if moving with the same speed the whole time.
func Travel(c Circle, travelVec Pt, travelLen Int, obstacles []Square) (newPos Pt, newTravelVec Pt, stop bool) {
	oldPos := c.Center

	bounces := 0
	for {
		// Given an original position and a travel vector, compute the new
		// position.
		newPos = oldPos.Plus(travelVec.Times(travelLen).DivBy(travelVec.Len()))

		// Check if we can travel to newPos without collision.
		// CircleSquareCollision doesn't return oldPos as a collision point.
		intersects, circlePositionAtCollision, collisionNormal :=
			CircleSquaresCollision(oldPos, newPos, c.Diameter, obstacles)
		if !intersects {
			// No collision, so we're fine, newPos is the final position.
			if bounces > 0 {
				fmt.Print(bounces, " ")
			}
			return newPos, travelVec, false
		}
		bounces++

		// We collided. We were supposed to travel travelLen but we only
		// travelled part of that then collided.
		// Find out how much we travelled.
		travelledLen := oldPos.To(circlePositionAtCollision).Len()
		// Move to the point where we collided.
		oldPos = circlePositionAtCollision
		// Move a little away from the collision point. If we ever let the ball
		// occupy a position where it is colliding with an obstacle, we get in
		// all sorts of trouble with edge cases. All we want is to move the ball
		// 1 integer unit away from the obstacle. We know what "away" means
		// because we have the collision normal.
		maxCoord := Max(collisionNormal.X.Abs(), collisionNormal.Y.Abs())
		offset := Pt{collisionNormal.X.DivBy(maxCoord), collisionNormal.Y.DivBy(maxCoord)}
		// Offset is now one of these: {1, 0}, {-1, 0}, {0, 1}, {0, -1}
		oldPos.Add(offset)

		// Update the travel length.
		travelLen.Subtract(travelledLen)
		// Update the travel direction.
		travelVec.Reflect(collisionNormal)
	}
}

func UpdateBallPositions(balls []Ball, s []Square, dec Int) {
	// update the state of each ball (move it, make it collectible)
	for idx := range balls {
		ball := &balls[idx]
		if ball.Speed.Gt(I(0)) {
			// move the ball
			var stop bool
			ball.Bounds.Center, ball.MoveDir, stop = Travel(ball.Bounds, ball.MoveDir, ball.Speed, s)

			if stop {
				ball.Speed = I(0)
			} else {
				// decrease speed by some deceleration
				ball.Speed.Subtract(dec)
				if ball.Speed.Lt(I(0)) {
					ball.Speed = I(0)
				}
			}

		}
		if !ball.CanBeCollected && ball.Speed.Lt(CU(10)) {
			ball.CanBeCollected = true
		}
	}
}

func MovePlayer(player *Player, newPos Pt, squares []Square) {
	oldPos := player.Bounds.Center

	intersects, circlePositionAtCollision, collisionNormal :=
		CircleSquaresCollision(oldPos, newPos, player.Bounds.Diameter, squares)

	if !intersects {
		player.Bounds.Center = newPos
	} else {
		adjustedNewPos := circlePositionAtCollision
		maxCoord := Max(collisionNormal.X.Abs(), collisionNormal.Y.Abs())
		offset := Pt{collisionNormal.X.DivBy(maxCoord), collisionNormal.Y.DivBy(maxCoord)}
		// Offset is now one of these: {1, 0}, {-1, 0}, {0, 1}, {0, -1}
		adjustedNewPos.Add(offset.Times(I(50)))

		// Find out if the point where we move to is outside the collision zone.
		// Unfortunately the only way to test the collision zone is with a travel
		// line. This means I have to test by travelling from oldPos to the new pos.
		//intersects2, circlePositionAtCollision2, collisionNormal2 :=
		//	CircleSquaresCollision(oldPos, adjustedNewPos, player.Bounds.Diameter, squares)
		//if intersects2 {
		//	fmt.Print("i ", circlePositionAtCollision2, collisionNormal2)
		//}

		player.Bounds.Center = adjustedNewPos
	}
}

func MoveStraightLine(start, end Pt) (input PlayerInput) {
	dx := end.X.Minus(start.X)
	dy := end.Y.Minus(start.Y)

	tol := U(2) // Should be greater than half of the player's speed.
	if dx.Lt(tol.Negative()) {
		input.MoveLeft = true
	} else if dx.Gt(tol) {
		input.MoveRight = true
	}
	if dy.Lt(tol.Negative()) {
		input.MoveUp = true
	} else if dy.Gt(tol) {
		input.MoveDown = true
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

func HandlePlayerInput(player *Player, balls *[]Ball, input PlayerInput,
	ballSpeed Int, squares []Square) {

	// Try horizontal movement first.
	newPosX := player.Bounds.Center
	if input.MoveRight {
		newPosX.X.Add(player.Speed)
	}
	if input.MoveLeft {
		newPosX.X.Subtract(player.Speed)
	}
	MovePlayer(player, newPosX, squares)

	// Now try vertical movement.
	newPosY := player.Bounds.Center
	if input.MoveUp {
		newPosY.Y.Subtract(player.Speed)
	}
	if input.MoveDown {
		newPosY.Y.Add(player.Speed)
	}
	MovePlayer(player, newPosY, squares)

	if input.Shoot {
		ShootBall(player, balls, input.ShootPt, ballSpeed)
	}
}

func PlayerAndBallAreTouching(player Player, ball Ball) bool {
	return CirclesIntersect(player.Bounds, ball.Bounds)
}

func FriendlyBall(player Player, ball Ball) bool {
	return player.BallType.Eq(ball.Type)
}

func HandlePlayerBallInteraction(player *Player, balls *[]Ball) {
	toBeDeleted := make([]bool, len(*balls))
	for idx, ball := range *balls {
		if !PlayerAndBallAreTouching(*player, ball) {
			continue
		}

		if FriendlyBall(*player, ball) {
			if ball.CanBeCollected {
				toBeDeleted[idx] = true
				// Disable this for debugging purposes.
				// player.NBalls.Inc()
			}
		} else {
			if player.Health.Gt(I(0)) {
				player.Health.Dec()
			}
			toBeDeleted[idx] = true
			// Disable this for debugging purposes.
			//player.NBalls.Inc()
		}
	}

	var newBalls []Ball
	for idx, ball := range *balls {
		if !toBeDeleted[idx] {
			newBalls = append(newBalls, ball)
		}
	}
	*balls = newBalls
}

func ObstaclesToSquares(obstacles Matrix, obstacleSize Int) (squares []Square) {
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

func (w *World) Step(input *Input, frameIdx int) {
	squares := ObstaclesToSquares(w.Obstacles, w.ObstacleSize)

	HandlePlayerInput(&w.Player1, &w.Balls, input.Player1Input, w.BallSpeed, squares)
	if frameIdx == 10 {
		//ShootBallDebug(&w.Balls, UPt(200, 250), UPt(1000, 2000), MU(200000))
	}
	HandlePlayerInput(&w.Player2, &w.Balls, input.Player2Input, w.BallSpeed, squares)

	UpdateBallPositions(w.Balls, squares, w.BallDec)
	HandlePlayerBallInteraction(&w.Player1, &w.Balls)
	HandlePlayerBallInteraction(&w.Player2, &w.Balls)
}
