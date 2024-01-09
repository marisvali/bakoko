package bakoko

import (
	"bytes"
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
}

type Matrix struct {
	cells []Int
	nRows Int
	nCols Int
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

func (m *Matrix) Set(x, y, val Int) {
	m.cells[y.Times(m.nCols).Plus(x).ToInt64()] = val
}

func (m *Matrix) Get(x, y Int) Int {
	return m.cells[y.Times(m.nCols).Plus(x).ToInt64()]
}

func (m *Matrix) NRows() Int {
	return m.nRows
}

func (m *Matrix) NCols() Int {
	return m.nCols
}

type World struct {
	Player1      Player
	Player2      Player
	Balls        []Ball
	Over         Int
	Obstacles    Matrix
	ObstacleSize Int
	Obs          []Square
}

type PlayerInput struct {
	MoveLeft  bool
	MoveRight bool
	MoveUp    bool
	MoveDown  bool
	Shoot     bool
	ShootPt   Pt
	Quit      bool
}

type Input struct {
	Player1Input PlayerInput
	Player2Input PlayerInput
}

func (w *World) Serialize() []byte {
	buf := new(bytes.Buffer)
	Serialize(buf, w.Player1)
	Serialize(buf, w.Player2)
	Serialize(buf, I(int64(len(w.Balls))))
	Serialize(buf, w.Balls)
	w.Obstacles.Serialize(buf)
	Serialize(buf, w.ObstacleSize)
	Serialize(buf, I(int64(len(w.Obs))))
	Serialize(buf, w.Obs)
	return buf.Bytes()
}

func (w *World) Deserialize(buf *bytes.Buffer) {
	Deserialize(buf, &w.Player1)
	Deserialize(buf, &w.Player2)
	var lenBalls Int
	Deserialize(buf, &lenBalls)
	w.Balls = make([]Ball, lenBalls.ToInt64())
	Deserialize(buf, w.Balls)
	w.Obstacles.Deserialize(buf)
	Deserialize(buf, &w.ObstacleSize)
	var lenObs Int
	Deserialize(buf, &lenObs)
	w.Obs = make([]Square, lenObs.ToInt64())
	Deserialize(buf, w.Obs)
}

func (w *World) SerializeToFile(filename string) {
	data := w.Serialize()
	WriteFile(filename, data)
}

func (w *World) DeserializeFromFile(filename string) {
	buf := bytes.NewBuffer(ReadFile(filename))
	w.Deserialize(buf)
}

func (i *Input) SerializeToFile(filename string) {
	buf := new(bytes.Buffer)
	Serialize(buf, i)
	WriteFile(filename, buf.Bytes())
}

func (i *Input) DeserializeFromFile(filename string) {
	data := ReadFile(filename)
	buf := bytes.NewBuffer(data)
	Deserialize(buf, i)
}

func ShootBall(player *Player, balls *[]Ball, pt Pt) {
	if player.NBalls.Leq(I(0)) {
		return
	}

	moveDir := player.Bounds.Center.To(pt)
	//moveDir.SetLen(MU(6000))
	moveDir.SetLen(U(1))
	speed := MU(6000)

	ball := Ball{
		//Center:            Pt{player.Center.X + (player.Diameter+30*Unit)/2 + 2*Unit, player.Center.Y},
		Bounds: Circle{
			Center:   player.Bounds.Center,
			Diameter: U(30)},
		MoveDir:        moveDir,
		Speed:          speed,
		CanBeCollected: false,
		Type:           player.BallType,
	}
	*balls = append(*balls, ball)
	player.NBalls.Dec()
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

func SerializeInputs(inputs []Input, filename string) {
	buf := new(bytes.Buffer)
	Serialize(buf, int64(len(inputs)))
	Serialize(buf, inputs)
	WriteFile(filename, buf.Bytes())
}

func DeserializeInputs(filename string) []Input {
	var inputs []Input
	buf := bytes.NewBuffer(ReadFile(filename))
	var lenInputs Int
	Deserialize(buf, &lenInputs)
	inputs = make([]Input, lenInputs.ToInt64())
	Deserialize(buf, inputs)
	return inputs
}

// Compute the new position of a circle travelling from its current position
// along travelVec to a new position. Collisions are handled so that the final
// position takes into account any obstacles hit along the way.
// The logic of this function is that the circle travels for a length of
// travelLen in total and has no concept of time. So you can say it treats
// the movement as uniform, as if moving with the same speed the whole time.
func Travel(c Circle, travelVec Pt, travelLen Int, obstacles []Square) (newPos Pt, newTravelVec Pt) {
	oldPos := c.Center

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
			return newPos, travelVec
		}

		// We collided. We were supposed to travel travelLen but we only
		// travelled part of that then collided.
		// Find out how much we travelled.
		travelledLen := oldPos.To(circlePositionAtCollision).Len()
		// Move to the point where we collided.
		oldPos = circlePositionAtCollision
		// Update the travel length.
		travelLen.Subtract(travelledLen)
		// Update the travel direction.
		travelVec.Reflect(collisionNormal)
	}
}

func UpdateBallPositions(balls []Ball, s []Square) {
	// update the state of each ball (move it, make it collectible)
	for idx := range balls {
		ball := &balls[idx]
		if ball.Speed.Gt(I(0)) {
			// move the ball
			ball.Bounds.Center, ball.MoveDir = Travel(ball.Bounds, ball.MoveDir, ball.Speed, s)

			// decrease speed by some deceleration
			ball.Speed.Subtract(MU(60))
			if ball.Speed.Lt(I(0)) {
				ball.Speed = I(0)
			}
		}
		if !ball.CanBeCollected && ball.Speed.Lt(MU(100)) {
			ball.CanBeCollected = true
		}
	}
}

func HandlePlayerInput(player *Player, balls *[]Ball, input PlayerInput) {
	if input.MoveRight {
		player.Bounds.Center.X.Add(U(3))
	}
	if input.MoveLeft {
		player.Bounds.Center.X.Subtract(U(3))
	}
	if input.MoveUp {
		player.Bounds.Center.Y.Subtract(U(3))
	}
	if input.MoveDown {
		player.Bounds.Center.Y.Add(U(3))
	}
	if input.Shoot {
		ShootBall(player, balls, input.ShootPt)
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
				player.NBalls.Inc()
			}
		} else {
			if player.Health.Gt(I(0)) {
				player.Health.Dec()
			}
			toBeDeleted[idx] = true
			player.NBalls.Inc()
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

func (w *World) Step(input *Input, frameIdx int) {
	HandlePlayerInput(&w.Player1, &w.Balls, input.Player1Input)
	if frameIdx == 10 {
		ShootBallDebug(&w.Balls, UPt(200, 250), UPt(200, 1000), MU(60000))
	}
	HandlePlayerInput(&w.Player2, &w.Balls, input.Player2Input)
	UpdateBallPositions(w.Balls, w.Obs)
	HandlePlayerBallInteraction(&w.Player1, &w.Balls)
	HandlePlayerBallInteraction(&w.Player2, &w.Balls)
}
