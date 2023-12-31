package bakoko

import (
	"bytes"
	. "playful-patterns.com/bakoko/ints"
)

type Ball struct {
	Pos            Pt
	Diameter       Int
	Speed          Pt
	CanBeCollected bool
}

type Player struct {
	Pos      Pt
	Diameter Int
	NBalls   Int
}

type World struct {
	Player1 Player
	Player2 Player
	Balls   []Ball
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

type Input struct {
	MoveLeft  bool
	MoveRight bool
	MoveUp    bool
	MoveDown  bool
	Shoot     bool
	ShootPt   Pt
	Quit      bool
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

func HandlePlayerBallInteraction(player *Player, balls *[]Ball) {
	var newBalls []Ball
	for _, ball := range *balls {
		maxDist := ball.Diameter.Plus(player.Diameter).DivBy(I(2))
		squaredMaxDist := maxDist.Sqr()
		if player.Pos.SquaredDistTo(ball.Pos).Gt(squaredMaxDist) || !ball.
			CanBeCollected {
			newBalls = append(newBalls, ball)
		} else {
			player.NBalls.Inc()
		}
	}
	*balls = newBalls
}

func ShootBall(player *Player, balls *[]Ball, pt Pt) {
	if player.NBalls.Leq(I(0)) {
		return
	}

	speed := player.Pos.To(pt)
	speed.SetLen(MU(6000))

	ball := Ball{
		//Pos:            Pt{player.Pos.X + (player.Diameter+30*Unit)/2 + 2*Unit, player.Pos.Y},
		Pos:            player.Pos,
		Diameter:       U(30),
		Speed:          speed,
		CanBeCollected: false,
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
	if input.MoveRight {
		w.Player1.Pos.X.Add(U(3))
	}
	if input.MoveLeft {
		w.Player1.Pos.X.Subtract(U(3))
	}
	if input.MoveUp {
		w.Player1.Pos.Y.Subtract(U(3))
	}
	if input.MoveDown {
		w.Player1.Pos.Y.Add(U(3))
	}
	if input.Shoot {
		ShootBall(&w.Player1, &w.Balls, input.ShootPt)
	}

	for idx := range w.Balls {
		ball := &w.Balls[idx]
		if ball.Speed.SquaredLen().Gt(I(0)) {
			ball.Pos.Add(ball.Speed)
			ball.Speed.AddLen(MU(-600))
		}
		if !ball.CanBeCollected && ball.Speed.SquaredLen().Lt(MU(100)) {
			ball.CanBeCollected = true
		}
	}

	HandlePlayerBallInteraction(&w.Player1, &w.Balls)
	HandlePlayerBallInteraction(&w.Player2, &w.Balls)
}
