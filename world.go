package bakoko

import (
	"bytes"
	"math"
)

type Real = float64
type Int = int64

const Unit = Int(1000)

type Point struct {
	X, Y Int
}

func (p *Point) DistTo(other Point) Int {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return Int(math.Sqrt(Real(dx*dx + dy*dy)))
}

func (p *Point) Add(other Point) {
	p.X += other.X
	p.Y += other.Y
}

func (p *Point) Mul(factor Real) Point {
	p.X = Int(Real(p.X) * factor)
	p.Y = Int(Real(p.Y) * factor)
	return *p
}

func (p *Point) Len() Int {
	return Int(math.Sqrt(Real(p.X*p.X + p.Y*p.Y)))
}

func (p *Point) To(other Point) Point {
	return Point{other.X - p.X, other.Y - p.Y}
}

func (p *Point) Norm() Point {
	p.Mul(Real(Unit) / Real(p.Len()))
	return *p
}

type Ball struct {
	Pos            Point
	Diameter       Int
	Speed          Point
	CanBeCollected bool
}

type Character struct {
	Pos      Point
	Diameter Int
	NBalls   Int
}

type World struct {
	Player1 Character
	Player2 Character
	Balls   []Ball
}

func (w *World) Serialize() []byte {
	buf := new(bytes.Buffer)
	Serialize(buf, w.Player1)
	Serialize(buf, w.Player2)
	Serialize(buf, Int(len(w.Balls)))
	Serialize(buf, w.Balls)
	return buf.Bytes()
}

func (w *World) Deserialize(buf *bytes.Buffer) {
	Deserialize(buf, &w.Player1)
	Deserialize(buf, &w.Player2)
	var lenBalls Int
	Deserialize(buf, &lenBalls)
	w.Balls = make([]Ball, lenBalls)
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
	ShootPt   Point
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

func HandlePlayerBallInteraction(player *Character, balls *[]Ball) {
	var newBalls []Ball
	for _, ball := range *balls {
		maxDist := (ball.Diameter + player.Diameter) / 2
		if player.Pos.DistTo(ball.Pos) > maxDist || !ball.CanBeCollected {
			newBalls = append(newBalls, ball)
		} else {
			player.NBalls++
		}
	}
	*balls = newBalls
}

func ShootBall(player *Character, balls *[]Ball, pt Point) {
	if player.NBalls <= 0 {
		return
	}

	speed := player.Pos.To(pt)
	speed.Norm()
	speed.Mul(3)

	ball := Ball{
		//Pos:            Point{player.Pos.X + (player.Diameter+30*Unit)/2 + 2*Unit, player.Pos.Y},
		Pos:            player.Pos,
		Diameter:       30 * Unit,
		Speed:          speed,
		CanBeCollected: false,
	}
	*balls = append(*balls, ball)
	player.NBalls--
}

func (w *World) Step(input *Input, frameIdx int) {
	if input.MoveRight {
		w.Player1.Pos.X += Unit
	}
	if input.MoveLeft {
		w.Player1.Pos.X -= Unit
	}
	if input.MoveUp {
		w.Player1.Pos.Y -= Unit
	}
	if input.MoveDown {
		w.Player1.Pos.Y += Unit
	}
	if input.Shoot || frameIdx == 20 {
		ShootBall(&w.Player1, &w.Balls, input.ShootPt.Mul(Real(Unit)))
	}

	for idx := range w.Balls {
		ball := &w.Balls[idx]
		if ball.Speed.Len() > 0 {
			ball.Pos.Add(ball.Speed)
			factor := Real(ball.Speed.Len()-30) / Real(ball.Speed.Len())
			ball.Speed.Mul(factor)
		}
		if !ball.CanBeCollected && ball.Speed.Len() < 100 {
			ball.CanBeCollected = true
		}
	}

	HandlePlayerBallInteraction(&w.Player1, &w.Balls)
	HandlePlayerBallInteraction(&w.Player2, &w.Balls)
}
