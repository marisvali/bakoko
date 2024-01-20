package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
	"io"
	"log"
	"net"
	"os"
	. "playful-patterns.com/bakoko"
	. "playful-patterns.com/bakoko/ints"
	"slices"
	"time"
)

func colorHex(hexVal int) color.Color {
	if hexVal < 0x000000 || hexVal > 0xFFFFFF {
		panic(fmt.Sprintf("Invalid HEX value for color: %d", hexVal))
	}
	r := uint8(hexVal & 0xFF0000 >> 16)
	g := uint8(hexVal & 0x00FF00 >> 8)
	b := uint8(hexVal & 0x0000FF)
	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
}

type simulationPeer struct {
	endpoint string
	conn     net.Conn
}

// Doesn't matter if this fails.
func (p *simulationPeer) getWorld(w *World) {
	// Don't do anything if we don't have a peer.
	// The communication between us and the peer is always that:
	// - we connect to the peer
	// - we send input to the peer
	// - we get an ouput from the peer
	// If the peer disconnects in middle of that, we start from the beginning,
	// we don't accept a connection then continue with getting the output.
	if p.conn == nil {
		return
	}

	data, err := ReadData(p.conn)
	// If there was an error, assume the peer is no longer available.
	// Invalidate the connection and try again later.
	if err != nil {
		p.conn = nil
		log.Println("lost connection")
		return
	}

	w.Deserialize(bytes.NewBuffer(data))
}

// Try to send an input to the peer, but don't block.
func (p *simulationPeer) sendInput(input *PlayerInput) {
	// If we don't have a peer, connect to one.
	if p.conn == nil {
		var err error
		p.conn, err = net.DialTimeout("tcp", p.endpoint, 5*time.Millisecond)

		// If connection took too long or failed, screw it.
		// We'll try again later.
		if err != nil {
			//log.Println("could not connect!")
			return
		}
	}
	//log.Println("connection established!")

	// We have a connection, try to send our input.
	buf := new(bytes.Buffer)
	Serialize(buf, input)

	err := WriteData(p.conn, buf.Bytes())
	// If there was an error, assume the peer is no longer available.
	// Invalidate the connection and try again later.
	if err != nil {
		p.conn = nil
		log.Println("lost connection")
	}
}

func (g *Game) Update() error {
	// Get keyboard input.
	var pressedKeys []ebiten.Key
	pressedKeys = inpututil.AppendPressedKeys(pressedKeys)
	//pressedKeys = inpututil.AppendJustPressedKeys(pressedKeys) //for debug purposes

	// Choose which is the active player based on Alt being pressed.
	playerInput := PlayerInput{}
	playerInput.MoveLeft = slices.Contains(pressedKeys, ebiten.KeyA)
	playerInput.MoveUp = slices.Contains(pressedKeys, ebiten.KeyW)
	playerInput.MoveDown = slices.Contains(pressedKeys, ebiten.KeyS)
	playerInput.MoveRight = slices.Contains(pressedKeys, ebiten.KeyD)

	var justPressedKeys []ebiten.Key
	justPressedKeys = inpututil.AppendJustPressedKeys(justPressedKeys)
	if slices.Contains(justPressedKeys, ebiten.KeyR) {
		g.loadGameData()
	}

	if g.gameDataChangedOnDisk() {
		g.loadGameData()
	}

	// Get mouse input.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		playerInput.Shoot = true
		x, y := ebiten.CursorPosition()
		// Translate from screen coordinates to in-world units.
		playerInput.ShootPt.X = g.ScreenToWorld(x)
		playerInput.ShootPt.Y = g.ScreenToWorld(y)
	}

	//g.w.Step(&input)
	g.peer.sendInput(&playerInput)
	//var w World
	//g.peer.getWorld(&w)
	g.peer.getWorld(&g.w)
	//input.SerializeToFile("input.bin")
	//TouchFile("input-ready")
	//WaitForFile("world-ready")
	//g.w.DeserializeFromFile("world.bin")
	return nil
}

//func DrawSprite(screen *ebiten.Image, img *ebiten.Image, pos Pt) {
//	op := &ebiten.DrawImageOptions{}
//	op.GeoM.Translate(pos.X.ToFloat64(), pos.Y.ToFloat64())
//	screen.DrawImage(img, op)
//}

func (g *Game) WorldToScreen(val Int) float64 {
	return val.ToFloat64() / Unit * g.data.ScaleFactor
}

func (g *Game) ScreenToWorld(val int) Int {
	return U(int64(float64(val) / g.data.ScaleFactor))
}

func (g *Game) DrawSprite(img *ebiten.Image,
	x float64, y float64, targetSize float64) {
	op := &ebiten.DrawImageOptions{}

	// Resize image to fit the targetSize of the circle we want to draw.
	// This kind of scaling is very useful during development when the final
	// sizes are not decided, and thus it's impossible to have final sprites.
	// For an actual release, scaling should be avoided.
	imgSize := img.Bounds().Size()
	newDx := targetSize / float64(imgSize.X)
	newDy := targetSize / float64(imgSize.Y)
	op.GeoM.Scale(newDx, newDy)

	// Place the image so that (x, y) falls at its center,
	// not its top-left corner.
	op.GeoM.Translate(x-targetSize/2, y-targetSize/2)

	g.screen.DrawImage(img, op)

	// Draw a small white rectangle in the center of the image,
	// to help debug issues with scaling and positioning.
	//dbgImg := ebiten.NewImage(3, 3)
	//dbgImg.Fill(color.RGBA{
	//	R: 255,
	//	G: 255,
	//	B: 255,
	//	A: 255,
	//})
	//op = &ebiten.DrawImageOptions{}
	//op.GeoM.Scale(newDx, newDy)
	//op.GeoM.Translate(x, y)
	//screen.DrawImage(dbgImg, op)
}

func (g *Game) DrawPlayer(
	playerImage *ebiten.Image,
	ballImage *ebiten.Image,
	healthImage *ebiten.Image,
	player *Player) {
	// Draw the player sprite.
	x := g.WorldToScreen(player.Bounds.Center.X)
	y := g.WorldToScreen(player.Bounds.Center.Y)
	diam := g.WorldToScreen(player.Bounds.Diameter)
	g.DrawSprite(playerImage, x, y, diam)

	// Draw a small sprite for each ball that the player has.
	for idx := int64(0); idx < player.NBalls.ToInt64(); idx++ {
		smallX := x - diam/2 - 10*g.data.ScaleFactor
		smallY := y + (float64(idx*12)+10)*g.data.ScaleFactor - diam/2
		smallDiam := float64(10) * g.data.ScaleFactor
		g.DrawSprite(ballImage, smallX, smallY, smallDiam)
	}

	// Draw a small sprite for each health point that the player has.
	for idx := int64(0); idx < player.Health.ToInt64(); idx++ {
		smallX := x + (float64(idx*12)+12)*g.data.ScaleFactor - diam/2
		smallY := y - diam/2 - 10*g.data.ScaleFactor
		smallDiam := float64(10) * g.data.ScaleFactor
		g.DrawSprite(healthImage, smallX, smallY, smallDiam)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.screen = screen

	// Background
	screen.Fill(g.background.At(0, 0))

	// Obstacle grid
	for y := I(0); y.Lt(g.w.Obstacles.NRows()); y.Inc() {
		for x := I(0); x.Lt(g.w.Obstacles.NCols()); x.Inc() {
			if g.w.Obstacles.Get(y, x).Eq(I(1)) {
				xScreen := g.WorldToScreen(x.Times(g.w.ObstacleSize).Plus(g.w.ObstacleSize.DivBy(I(2))))
				yScreen := g.WorldToScreen(y.Times(g.w.ObstacleSize).Plus(g.w.ObstacleSize.DivBy(I(2))))
				diameter := g.WorldToScreen(g.w.ObstacleSize)
				g.DrawSprite(g.obstacle, xScreen, yScreen, diameter)
			}
		}
	}

	// debug squares
	for _, sq := range g.w.Obs {
		xScreen := g.WorldToScreen(sq.Center.X)
		yScreen := g.WorldToScreen(sq.Center.Y)
		diameter := g.WorldToScreen(sq.Size)
		g.DrawSprite(g.obstacle, xScreen, yScreen, diameter)
	}

	// Players
	g.DrawPlayer(g.player1, g.ball1, g.health, &g.w.Player1)
	g.DrawPlayer(g.player2, g.ball2, g.health, &g.w.Player2)

	// Balls
	for _, ball := range g.w.Balls {
		ballImage := g.ball1
		if ball.Type.Eq(I(2)) {
			ballImage = g.ball2
		}
		g.DrawSprite(ballImage,
			g.WorldToScreen(ball.Bounds.Center.X),
			g.WorldToScreen(ball.Bounds.Center.Y),
			g.WorldToScreen(ball.Bounds.Diameter))
	}

	//img1 := ebiten.NewImage(50, 50)
	//img1.Fill(colorPrimary)
	//op := &ebiten.DrawImageOptions{}
	//op.GeoM.Translate(Real(g.w.Player1.Center.X), Real(g.w.Player1.Center.Y))
	//screen.DrawImage(img1, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("ActualTPS: %f", ebiten.ActualTPS()))

	g.screen = nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func init() {
}

type GameData struct {
	ScaleFactor  float64
	WindowWidth  int
	WindowHeight int
}

type Game struct {
	w          World
	peer       simulationPeer
	player1    *ebiten.Image
	player2    *ebiten.Image
	ball1      *ebiten.Image
	ball2      *ebiten.Image
	health     *ebiten.Image
	obstacle   *ebiten.Image
	background *ebiten.Image
	screen     *ebiten.Image
	data       GameData
	times      []time.Time
}

func loadImage(str string) *ebiten.Image {
	file, err := os.Open(str)
	defer file.Close()
	Check(err)

	img, _, err := image.Decode(file)
	Check(err)
	return ebiten.NewImageFromImage(img)
}

func loadJSON(filename string, v any) {
	file, err := os.Open(filename)
	Check(err)
	bytes, err := io.ReadAll(file)
	Check(err)
	err = json.Unmarshal(bytes, v)
	Check(err)
}

func (g *Game) gameDataChangedOnDisk() bool {
	files, err := os.ReadDir("data")
	Check(err)
	if len(files) != len(g.times) {
		g.times = make([]time.Time, len(files))
	}
	changed := false
	for idx, file := range files {
		info, err := file.Info()
		Check(err)
		if g.times[idx] != info.ModTime() {
			changed = true
			g.times[idx] = info.ModTime()
		}
	}
	return changed
}

func (g *Game) loadGameData() {
	g.ball1 = loadImage("data/ball1.png")
	g.ball2 = loadImage("data/ball2.png")
	g.player1 = loadImage("data/player1.png")
	g.player2 = loadImage("data/player2.png")
	g.health = loadImage("data/health.png")
	g.obstacle = loadImage("data/obstacle.png")
	g.background = loadImage("data/background.png")
	loadJSON("data/gui.json", &g.data)

	ebiten.SetWindowSize(g.data.WindowWidth, g.data.WindowHeight)
	ebiten.SetWindowTitle("Viewer")
	ebiten.SetWindowPosition(10, 1080-10-g.data.WindowHeight)
}

func main() {
	var g Game
	g.peer.endpoint = os.Args[1] // localhost:56901 or localhost:56902
	g.loadGameData()

	err := ebiten.RunGame(&g)
	Check(err)
}
