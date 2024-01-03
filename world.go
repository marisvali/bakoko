package bakoko

import (
	"bytes"
	. "playful-patterns.com/bakoko/ints"
)

type Ball struct {
	Type           Int
	Pos            Pt
	Diameter       Int
	Speed          Pt
	CanBeCollected bool
}

type Player struct {
	Pos      Pt
	Diameter Int
	NBalls   Int
	BallType Int
	Health   Int
}

type World struct {
	Player1 Player
	Player2 Player
	Balls   []Ball
	Over    Int
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
	return buf.Bytes()
}

func (w *World) Deserialize(buf *bytes.Buffer) {
	Deserialize(buf, &w.Player1)
	Deserialize(buf, &w.Player2)
	var lenBalls Int
	Deserialize(buf, &lenBalls)
	w.Balls = make([]Ball, lenBalls.ToInt64())
	Deserialize(buf, w.Balls)
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

	speed := player.Pos.To(pt)
	//speed.SetLen(MU(6000))
	speed.SetLen(MU(6000))

	ball := Ball{
		//Pos:            Pt{player.Pos.X + (player.Diameter+30*Unit)/2 + 2*Unit, player.Pos.Y},
		Pos:            player.Pos,
		Diameter:       U(30),
		Speed:          speed,
		CanBeCollected: false,
		Type:           player.BallType,
	}
	*balls = append(*balls, ball)
	player.NBalls.Dec()
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

func (w *World) Step(input *Input, frameIdx int) {
	HandlePlayerInput(&w.Player1, &w.Balls, input.Player1Input)
	HandlePlayerInput(&w.Player2, &w.Balls, input.Player2Input)
	UpdateBallPositions(w.Balls)
	HandlePlayerBallInteraction(&w.Player1, &w.Balls)
	HandlePlayerBallInteraction(&w.Player2, &w.Balls)
}

func UpdateBallPositions(balls []Ball) {
	// update the state of each ball (move it, make it collectible)
	for idx := range balls {
		ball := &balls[idx]
		if ball.Speed.SquaredLen().Gt(I(0)) {
			ball.Pos.Add(ball.Speed)
			ball.Speed.AddLen(MU(-60))
		}
		if !ball.CanBeCollected && ball.Speed.SquaredLen().Lt(MU(100)) {
			ball.CanBeCollected = true
		}
	}
}

func HandlePlayerInput(player *Player, balls *[]Ball, input PlayerInput) {
	if input.MoveRight {
		player.Pos.X.Add(U(3))
	}
	if input.MoveLeft {
		player.Pos.X.Subtract(U(3))
	}
	if input.MoveUp {
		player.Pos.Y.Subtract(U(3))
	}
	if input.MoveDown {
		player.Pos.Y.Add(U(3))
	}
	if input.Shoot {
		ShootBall(player, balls, input.ShootPt)
	}
}

func PlayerAndBallAreTouching(player Player, ball Ball) bool {
	maxDist := ball.Diameter.Plus(player.Diameter).DivBy(I(2))
	squaredMaxDist := maxDist.Sqr()
	return player.Pos.SquaredDistTo(ball.Pos).Leq(squaredMaxDist)
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
