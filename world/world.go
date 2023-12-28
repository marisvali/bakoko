package world

import (
	"bytes"
	"math"
	"playful-patterns.com/bakoko/utils"
)

type Point struct {
	X, Y int64
}

func (p *Point) DistTo(other Point) int64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return int64(math.Sqrt(float64(dx*dx + dy*dy)))
}

func (p *Point) Add(other Point) {
	p.X += other.X
	p.Y += other.Y
}

func (p *Point) Mul(factor float64) Point {
	p.X = int64(float64(p.X) * factor)
	p.Y = int64(float64(p.Y) * factor)
	return *p
}

func (p *Point) Len() int64 {
	return int64(math.Sqrt(float64(p.X*p.X + p.Y*p.Y)))
}

func (p *Point) To(other Point) Point {
	return Point{other.X - p.X, other.Y - p.Y}
}

func (p *Point) Norm() Point {
	p.Mul(float64(1000) / float64(p.Len()))
	return *p
}

type Ball struct {
	Pos            Point
	Diameter       int64
	Speed          Point
	CanBeCollected bool
}

type Character struct {
	Pos      Point
	Diameter int64
	NBalls   int64
}

type World struct {
	Player1 Character
	Player2 Character
	Balls   []Ball
}

func (w *World) Serialize() []byte {
	buf := new(bytes.Buffer)
	utils.Serialize(buf, w.Player1)
	utils.Serialize(buf, w.Player2)
	utils.Serialize(buf, int64(len(w.Balls)))
	utils.Serialize(buf, w.Balls)
	return buf.Bytes()
}

func (w *World) Deserialize(buf *bytes.Buffer) {
	utils.Deserialize(buf, &w.Player1)
	utils.Deserialize(buf, &w.Player2)
	var lenBalls int64
	utils.Deserialize(buf, &lenBalls)
	w.Balls = make([]Ball, lenBalls)
	utils.Deserialize(buf, w.Balls)
}

func (w *World) SerializeToFile(filename string) {
	data := w.Serialize()
	utils.WriteFile(filename, data)
}

func (w *World) DeserializeFromFile(filename string) {
	buf := bytes.NewBuffer(utils.ReadFile(filename))
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
	utils.Serialize(buf, i)
	utils.WriteFile(filename, buf.Bytes())
}

func (i *Input) DeserializeFromFile(filename string) {
	data := utils.ReadFile(filename)
	buf := bytes.NewBuffer(data)
	utils.Deserialize(buf, i)
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
		//Pos:            Point{player.Pos.X + (player.Diameter+30*1000)/2 + 2*1000, player.Pos.Y},
		Pos:            player.Pos,
		Diameter:       30 * 1000,
		Speed:          speed,
		CanBeCollected: false,
	}
	*balls = append(*balls, ball)
	player.NBalls--
}

func (w *World) Step(input *Input, frameIdx int) error {
	if input.MoveRight {
		w.Player1.Pos.X += 1000
	}
	if input.MoveLeft {
		w.Player1.Pos.X -= 1000
	}
	if input.MoveUp {
		w.Player1.Pos.Y -= 1000
	}
	if input.MoveDown {
		w.Player1.Pos.Y += 1000
	}
	if input.Shoot || frameIdx == 20 {
		ShootBall(&w.Player1, &w.Balls, input.ShootPt.Mul(1000))
	}

	for idx := range w.Balls {
		ball := &w.Balls[idx]
		if ball.Speed.Len() > 0 {
			ball.Pos.Add(ball.Speed)
			factor := float64(ball.Speed.Len()-30) / float64(ball.Speed.Len())
			ball.Speed.Mul(factor)
		}
		if !ball.CanBeCollected && ball.Speed.Len() < 100 {
			ball.CanBeCollected = true
		}
	}

	HandlePlayerBallInteraction(&w.Player1, &w.Balls)
	HandlePlayerBallInteraction(&w.Player2, &w.Balls)
	return nil
}
