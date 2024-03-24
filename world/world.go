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
	Bounds            Circle
	NBalls            Int
	BallType          Int
	Health            Int
	Speed             Int
	State             Int
	StunnedImobilizes bool
	StunnedTime       Int
}

var PlayerRegular = I(0)
var PlayerStunned = I(1)

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
	BallDiameter Int
	DebugInfo    DebugInfo
	JustReloaded Int
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
	Pause     bool
}

func SerializeInputs(inputs []PlayerInput, filename string) {
	buf := new(bytes.Buffer)
	Serialize(buf, int64(len(inputs)))
	Serialize(buf, inputs)
	WriteFile(filename, buf.Bytes())
}

func DeserializeInputs(filename string) []PlayerInput {
	var inputs []PlayerInput
	buf := bytes.NewBuffer(ReadFile(filename))
	var lenInputs Int
	Deserialize(buf, &lenInputs)
	inputs = make([]PlayerInput, lenInputs.ToInt64())
	Deserialize(buf, inputs)
	return inputs
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
	Serialize(buf, w.JustReloaded)
	return buf.Bytes()
}

func (w *World) Deserialize(buf *bytes.Buffer) {
	Deserialize(buf, &w.Player1)
	Deserialize(buf, &w.Player2)
	DeserializeSlice(buf, &w.Balls)
	w.Obstacles.Deserialize(buf)
	Deserialize(buf, &w.ObstacleSize)
	Deserialize(buf, &w.JustReloaded)
}

func (w *World) ShootBall(player *Player, pt Pt) {
	if player.NBalls.Leq(I(0)) {
		return
	}

	moveDir := player.Bounds.Center.To(pt)
	//moveDir.SetLen(MU(6000))
	moveDir.SetLen(U(1))
	speed := w.BallSpeed

	// If the player moves right and shoots a ball to the right, the speed should compound.
	// How do I make this true?
	// I can only mess with the speed. The direction stays the same.
	// So, if I move right and I shoot the ball to the right,

	ball := Ball{
		//Pos:            Pt{player.Pos.X + (player.Diameter+30*Unit)/2 + 2*Unit, player.Pos.Y},
		Bounds: Circle{
			Center:   player.Bounds.Center,
			Diameter: w.BallDiameter},
		MoveDir:        moveDir,
		Speed:          speed,
		CanBeCollected: false,
		Type:           player.BallType,
	}
	w.Balls = append(w.Balls, ball)
	// Infinite balls, for debugging purposes.
	player.NBalls.Dec()
}

func ShootBallDebug(balls *[]Ball, orig, dest Pt, speed Int, ballDiameter Int) {
	moveDir := orig.To(dest)
	moveDir.SetLen(U(1))

	ball := Ball{
		Bounds: Circle{
			Center:   orig,
			Diameter: ballDiameter},
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
func (w *World) Travel(c Circle, travelVec Pt, travelLen Int) (newPos Pt, newTravelVec Pt, stop bool) {
	oldPos := c.Center

	bounces := 0
	for {
		// Given an original position and a travel vector, compute the new
		// position.
		newPos = oldPos.Plus(travelVec.Times(travelLen).DivBy(travelVec.Len()))
		obstacles := w.GetRelevantSquares(c.Diameter, oldPos, newPos)

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

func (w *World) UpdateBallPositions(balls []Ball, dec Int) {
	// update the state of each ball (move it, make it collectible)
	for idx := range balls {
		ball := &balls[idx]
		if ball.Speed.Gt(I(0)) {
			// move the ball
			var stop bool
			ball.Bounds.Center, ball.MoveDir, stop = w.Travel(ball.Bounds, ball.MoveDir, ball.Speed)

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

func (w *World) MovePlayer(player *Player, newPos Pt) {
	oldPos := player.Bounds.Center

	squares := w.GetRelevantSquares(player.Bounds.Diameter, oldPos, newPos)

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

func (w *World) HandlePlayerInput(player *Player, input PlayerInput) {

	if player.State.Eq(PlayerStunned) && player.StunnedImobilizes {
		return // Can't move or shoot while stunned.
	}

	// Try horizontal movement first.
	newPosX := player.Bounds.Center
	if input.MoveRight {
		newPosX.X.Add(player.Speed)
	}
	if input.MoveLeft {
		newPosX.X.Subtract(player.Speed)
	}
	w.MovePlayer(player, newPosX)

	// Now try vertical movement.
	newPosY := player.Bounds.Center
	if input.MoveUp {
		newPosY.Y.Subtract(player.Speed)
	}
	if input.MoveDown {
		newPosY.Y.Add(player.Speed)
	}
	w.MovePlayer(player, newPosY)

	if input.Shoot {
		w.ShootBall(player, input.ShootPt)
	}
}

func PlayerAndBallAreTouching(player Player, ball Ball) bool {
	return CirclesIntersect(player.Bounds, ball.Bounds)
}

func FriendlyBall(player Player, ball Ball) bool {
	return player.BallType.Eq(ball.Type)
}

func (w *World) HandlePlayerBallInteraction(player *Player, balls *[]Ball) {
	toBeDeleted := make([]bool, len(*balls))
	for idx, ball := range *balls {
		if !PlayerAndBallAreTouching(*player, ball) {
			continue
		}

		if FriendlyBall(*player, ball) {
			if ball.CanBeCollected {
				toBeDeleted[idx] = true
				// Disable this for debugging purposes.
				player.NBalls.Inc()
			}
		} else {
			if player.Health.Gt(I(0)) {
				player.Health.Dec()
				player.StunnedTime = I(30)
				player.State = PlayerStunned
			}
			toBeDeleted[idx] = true
			// Disable this for debugging purposes.
			player.NBalls.Inc()
		}
	}

	if player.StunnedTime.Gt(ZERO) {
		player.StunnedTime.Dec()
		if player.StunnedTime.Eq(ZERO) {
			player.State = PlayerRegular
		}
	}

	var newBalls []Ball
	for idx, ball := range *balls {
		if !toBeDeleted[idx] {
			newBalls = append(newBalls, ball)
		}
	}
	*balls = newBalls
	return
}

func (w *World) Step(input *Input, frameIdx int) {
	w.DebugInfo = DebugInfo{} // reset

	w.HandlePlayerInput(&w.Player1, input.Player1Input)
	if frameIdx == 10 {
		//ShootBallDebug(&w.Balls, UPt(200, 250), UPt(1000, 2000), MU(200000))
	}
	w.HandlePlayerInput(&w.Player2, input.Player2Input)

	w.UpdateBallPositions(w.Balls, w.BallDec)
	w.HandlePlayerBallInteraction(&w.Player1, &w.Balls)
	w.HandlePlayerBallInteraction(&w.Player2, &w.Balls)
}

func LoadWorld(w *World) {
	*w = World{} // Reset everything.

	data := loadWorldData(Home("world-data"))

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
	levelString := ReadAllText(Home(data.Level))
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

func (w *World) GetRelevantSquares(diameter Int, oldPos Pt, newPos Pt) (squares []Square) {
	// Compute the rectangle in which the travel will take place.
	x1, x2 := MinMax(oldPos.X, newPos.X)
	y1, y2 := MinMax(oldPos.Y, newPos.Y)

	// Expand the rectangle by half the diameter to make sure we get everything that is touched.
	radius := diameter.DivBy(TWO)
	// Expand the radius by (10% + 10) just to be sure we don't screw up because
	// of any tolerances. This function is used to get squares that MIGHT be
	// relevant. It's ok to get too many squares. It's not at all ok to get too
	// few.
	radius = radius.Plus(radius.DivBy(I(10))).Plus(I(10))
	x1.Subtract(radius)
	x2.Add(radius)
	y1.Subtract(radius)
	y2.Add(radius)

	// Convert the points to obstacle indexes.
	x1 = x1.DivBy(w.ObstacleSize)
	x2 = x2.DivBy(w.ObstacleSize)
	y1 = y1.DivBy(w.ObstacleSize)
	y2 = y2.DivBy(w.ObstacleSize)

	// Convert obstacles to squares.
	for row := y1; row.Leq(y2); row.Inc() {
		for col := x1; col.Leq(x2); col.Inc() {
			if w.Obstacles.Get(row, col).Neq(I(0)) {
				half := w.ObstacleSize.DivBy(I(2))
				square := Square{
					Center: Pt{
						col.Times(w.ObstacleSize).Plus(half),
						row.Times(w.ObstacleSize).Plus(half)},
					Size: w.ObstacleSize,
				}
				squares = append(squares, square)

				//var ds DebugSquare
				//ds.Square = square
				//ds.Col = color.RGBA{255, 0, 0, 255}
				//w.DebugInfo.Squares = append(w.DebugInfo.Squares, ds)
			}
		}
	}

	return

	//for y := I(0); y.Lt(w.Obstacles.NRows()); y.Inc() {
	//	for x := I(0); x.Lt(w.Obstacles.NCols()); x.Inc() {
	//		var pt DebugPoint
	//		pt.Pos = Pt{x.Times(w.ObstacleSize), y.Times(w.ObstacleSize)}
	//		offset := Pt{w.ObstacleSize.DivBy(TWO), w.ObstacleSize.DivBy(TWO)}
	//		pt.Pos.Add(offset)
	//		pt.Size = U(5)
	//		pt.Col = color.RGBA{0, 0, 255, 255}
	//		w.DebugInfo.Points = append(w.DebugInfo.Points, pt)
	//	}
	//}
	//
	//for y := y1; y.Leq(y2); y.Inc() {
	//	for x := x1; x.Leq(x2); x.Inc() {
	//		var pt DebugPoint
	//		pt.Pos = Pt{x.Times(w.ObstacleSize), y.Times(w.ObstacleSize)}
	//		offset := Pt{w.ObstacleSize.DivBy(TWO), w.ObstacleSize.DivBy(TWO)}
	//		pt.Pos.Add(offset)
	//		pt.Size = U(5)
	//		pt.Col = color.RGBA{255, 0, 255, 255}
	//		w.DebugInfo.Points = append(w.DebugInfo.Points, pt)
	//	}
	//}

	//w.DebugInfo.Points = append(w.DebugInfo.Points, DebugPoint{
	//	Pos:  Pt{x1, y1},
	//	Size: U(10),
	//	Col:  color.RGBA{255, 0, 0, 255},
	//})
	//w.DebugInfo.Points = append(w.DebugInfo.Points, DebugPoint{
	//	Pos:  Pt{x2, y2},
	//	Size: U(10),
	//	Col:  color.RGBA{255, 0, 0, 255},
	//})

	//x1 := oldPos.X.DivBy(obstacleSize)
	//y1 := oldPos.Y.DivBy(obstacleSize)
	//x2 := newPos.X.DivBy(obstacleSize)
	//y2 := newPos.Y.DivBy(obstacleSize)
	//
	//minX, maxX := MinMax(x1, x2)
	//minY, maxY := MinMax(y1, y2)
	//minX.Dec()
	//maxX.Add(TWO)
	//minY.Dec()
	//maxY.Add(TWO)
	//if minX.Lt(ZERO) {
	//	minX = ZERO
	//}
	//if minY.Lt(ZERO) {
	//	minY = ZERO
	//}
	//if maxX.Geq(obstacles.NCols()) {
	//	maxX = obstacles.NCols().Minus(ONE)
	//}
	//if maxY.Geq(obstacles.NRows()) {
	//	maxY = obstacles.NRows().Minus(ONE)
	//}
	//
	//half := obstacleSize.DivBy(I(2))
	//for row := minY; row.Leq(maxY); row.Inc() {
	//	for col := minX; col.Leq(maxX); col.Inc() {
	//		if obstacles.Get(row, col).Neq(ZERO) {
	//			squares = append(squares, Square{
	//				Center: Pt{col.Times(obstacleSize).Plus(half), row.Times(obstacleSize).Plus(half)},
	//				Size:   obstacleSize,
	//			})
	//		}
	//	}
	//}

	//half := obstacleSize.DivBy(I(2))
	//for row := ZERO; row.Lt(obstacles.NRows()); row.Inc() {
	//	for col := ZERO; col.Lt(obstacles.NCols()); col.Inc() {
	//		if obstacles.Get(row, col).Neq(ZERO) {
	//			squares = append(squares, Square{
	//				Center: Pt{col.Times(obstacleSize).Plus(half), row.Times(obstacleSize).Plus(half)},
	//				Size:   obstacleSize,
	//			})
	//		}
	//	}
	//}
	return
}
