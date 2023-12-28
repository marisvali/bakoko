package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
	"os"
	"playful-patterns.com/bakoko/utils"
	"playful-patterns.com/bakoko/world"
	"slices"
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

var colorPrimary = colorHex(0x05B2DC)
var colorPrimaryDark1 = colorHex(0x026d88)
var colorPrimaryDark2 = colorHex(0x002f3c)
var colorPrimaryLight1 = colorHex(0x76cae7)
var colorPrimaryLight2 = colorHex(0xb4e1f2)

var colorSecondary = colorHex(0xf52d00)
var colorSecondaryDark1 = colorHex(0x981800)
var colorSecondaryDark2 = colorHex(0x440600)
var colorSecondaryLight1 = colorHex(0xff7d64)
var colorSecondaryLight2 = colorHex(0xffb7a7)

var colorNeutral = colorHex(0x191308)
var colorNeutralLight1 = colorHex(0x2e2e2e)
var colorNeutralLight2 = colorHex(0x808080)
var colorNeutralLight3 = colorHex(0xdedede)

func (g *Game) Update() error {
	var input world.Input

	// Get keyboard input.
	var keys []ebiten.Key
	keys = inpututil.AppendPressedKeys(keys)
	input.MoveLeft = slices.Contains(keys, ebiten.KeyA)
	input.MoveUp = slices.Contains(keys, ebiten.KeyW)
	input.MoveDown = slices.Contains(keys, ebiten.KeyS)
	input.MoveRight = slices.Contains(keys, ebiten.KeyD)

	var keys2 []ebiten.Key
	keys2 = inpututil.AppendJustPressedKeys(keys2)

	// Get mouse input.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		input.Shoot = true
		x, y := ebiten.CursorPosition()
		input.ShootPt.X = int64(x)
		input.ShootPt.Y = int64(y)
	}

	//g.w.Step(&input)
	input.SerializeToFile("input.bin")
	utils.TouchFile("input-ready")
	utils.WaitForFile("world-ready")
	g.w.DeserializeFromFile("world.bin")
	return nil
}

func DrawSprite(screen *ebiten.Image, img *ebiten.Image, pos world.Point) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(pos.X), float64(pos.Y))
	screen.DrawImage(img, op)
}

func DrawCircle(screen *ebiten.Image, img *ebiten.Image, pos world.Point, radius float64) {
	op := &ebiten.DrawImageOptions{}
	size := img.Bounds().Size()
	newDx := radius / float64(1000) / float64(size.X)
	newDy := radius / float64(1000) / float64(size.Y)
	op.GeoM.Scale(newDx, newDy)
	// Have the pos indicate the center, not the top-left.
	//op.GeoM.Translate(float64(pos.X)-newDx/2, float64(pos.Y)-newDy/2)
	op.GeoM.Translate((float64(pos.X)-radius/2)/float64(1000), (float64(pos.Y)-radius/2)/float64(1000))
	screen.DrawImage(img, op)

	dbgImg := ebiten.NewImage(3, 3)
	dbgImg.Fill(color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	})
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(newDx, newDy)
	op.GeoM.Translate(float64(pos.X), float64(pos.Y))
	screen.DrawImage(dbgImg, op)
}

func DrawPlayer(screen *ebiten.Image, playerImage *ebiten.Image, player *world.Character) {
	DrawCircle(screen, playerImage, player.Pos, float64(player.Diameter))
	for idx := int64(0); idx < player.NBalls; idx++ {
		DrawCircle(screen, playerImage,
			world.Point{X: player.Pos.X - player.Diameter/2 - 10*1000, Y: player.Pos.Y + idx*12*1000 - player.Diameter/2 + 10*1000}, float64(10*1000))
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background.
	screen.Fill(colorNeutralLight1)

	DrawPlayer(screen, g.player1, &g.w.Player1)
	DrawPlayer(screen, g.player2, &g.w.Player2)
	for _, ball := range g.w.Balls {
		DrawCircle(screen, g.ball, ball.Pos, float64(ball.Diameter))
	}
	//img1 := ebiten.NewImage(50, 50)
	//img1.Fill(colorPrimary)
	//op := &ebiten.DrawImageOptions{}
	//op.GeoM.Translate(float64(g.w.Player1.Pos.X), float64(g.w.Player1.Pos.Y))
	//screen.DrawImage(img1, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("ActualTPS: %f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func init() {
}

type Game struct {
	w       world.World
	player1 *ebiten.Image
	player2 *ebiten.Image
	ball    *ebiten.Image
}

func loadImage(str string) *ebiten.Image {
	file, err := os.Open(str)
	defer file.Close()
	utils.Check(err)

	img, _, err := image.Decode(file)
	utils.Check(err)
	return ebiten.NewImageFromImage(img)
}

func main() {
	var g Game
	g.ball = loadImage("ball.png")
	g.player1 = loadImage("player1.png")
	g.player2 = loadImage("player2.png")
	ebiten.SetWindowSize(300, 300)
	ebiten.SetWindowTitle("Viewer")
	err := ebiten.RunGame(&g)
	utils.Check(err)
}
