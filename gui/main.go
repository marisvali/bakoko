package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
	"os"
	. "playful-patterns.com/bakoko"
	. "playful-patterns.com/bakoko/ints"
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
	var input Input

	// Get keyboard input.
	var pressedKeys []ebiten.Key
	pressedKeys = inpututil.AppendPressedKeys(pressedKeys)
	// Choose which is the active player based on Alt being pressed.
	playerInput := &input.Player1Input
	if slices.Contains(pressedKeys, ebiten.KeyAlt) {
		playerInput = &input.Player2Input
	}

	playerInput.MoveLeft = slices.Contains(pressedKeys, ebiten.KeyA)
	playerInput.MoveUp = slices.Contains(pressedKeys, ebiten.KeyW)
	playerInput.MoveDown = slices.Contains(pressedKeys, ebiten.KeyS)
	playerInput.MoveRight = slices.Contains(pressedKeys, ebiten.KeyD)

	var justPressedKeys []ebiten.Key
	justPressedKeys = inpututil.AppendJustPressedKeys(justPressedKeys)

	// Get mouse input.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		playerInput.Shoot = true
		x, y := ebiten.CursorPosition()
		// Translate from screen coordinates to in-world units.
		playerInput.ShootPt.X = ScreenToWorld(x)
		playerInput.ShootPt.Y = ScreenToWorld(y)
	}

	//g.w.Step(&input)
	input.SerializeToFile("input.bin")
	TouchFile("input-ready")
	WaitForFile("world-ready")
	g.w.DeserializeFromFile("world.bin")
	return nil
}

func DrawSprite(screen *ebiten.Image, img *ebiten.Image, pos Pt) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X.ToFloat64(), pos.Y.ToFloat64())
	screen.DrawImage(img, op)
}

func WorldToScreen(val Int) int {
	return int(val.ToInt64()) / Unit
}

func WorldToScreenFloat(val Int) float64 {
	return val.ToFloat64() / Unit
}

func ScreenToWorld(val int) Int {
	return U(int64(val))
}

func DrawCircle(screen *ebiten.Image, img *ebiten.Image, x float64, y float64,
	diameter float64) {
	op := &ebiten.DrawImageOptions{}

	// Resize image to fit the diameter of the circle we want to draw.
	// This kind of scaling is very useful during development when the final
	// sizes are not decided, and thus it's impossible to have final sprites.
	// For an actual release, scaling should be avoided.
	size := img.Bounds().Size()
	newDx := diameter / float64(size.X)
	newDy := diameter / float64(size.Y)
	op.GeoM.Scale(newDx, newDy)

	// Place the image so that (x, y) falls at its center,
	// not its top-left corner.
	op.GeoM.Translate(x-diameter/2, y-diameter/2)

	screen.DrawImage(img, op)

	// Draw a small white rectangle in the center of the image,
	// to help debug issues with scaling and positioning.
	dbgImg := ebiten.NewImage(3, 3)
	dbgImg.Fill(color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	})
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(newDx, newDy)
	op.GeoM.Translate(x, y)
	screen.DrawImage(dbgImg, op)
}

func DrawPlayer(
	screen *ebiten.Image,
	playerImage *ebiten.Image,
	ballImage *ebiten.Image,
	healthImage *ebiten.Image,
	player *Player) {
	// Draw the player sprite.
	x := WorldToScreenFloat(player.Pos.X)
	y := WorldToScreenFloat(player.Pos.Y)
	diam := WorldToScreenFloat(player.Diameter)
	DrawCircle(screen, playerImage, x, y, diam)

	// Draw a small sprite for each ball that the player has.
	for idx := int64(0); idx < player.NBalls.ToInt64(); idx++ {
		smallX := x - diam/2 - 10
		smallY := y + float64(idx*12) - diam/2 + 10
		smallDiam := float64(10)
		DrawCircle(screen, ballImage, smallX, smallY, smallDiam)
	}

	// Draw a small sprite for each health point that the player has.
	for idx := int64(0); idx < player.Health.ToInt64(); idx++ {
		smallX := x + float64(idx*12) - diam/2 + 15
		smallY := y - diam/2 - 10
		smallDiam := float64(10)
		DrawCircle(screen, healthImage, smallX, smallY, smallDiam)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background.
	screen.Fill(colorNeutralLight1)

	DrawPlayer(screen, g.player1, g.ball1, g.health, &g.w.Player1)
	DrawPlayer(screen, g.player2, g.ball1, g.health, &g.w.Player2)
	for _, ball := range g.w.Balls {
		ballImage := g.ball1
		if ball.Type.Eq(I(2)) {
			ballImage = g.ball2
		}
		DrawCircle(screen, ballImage,
			WorldToScreenFloat(ball.Pos.X),
			WorldToScreenFloat(ball.Pos.Y),
			WorldToScreenFloat(ball.Diameter))
	}
	//img1 := ebiten.NewImage(50, 50)
	//img1.Fill(colorPrimary)
	//op := &ebiten.DrawImageOptions{}
	//op.GeoM.Translate(Real(g.w.Player1.Pos.X), Real(g.w.Player1.Pos.Y))
	//screen.DrawImage(img1, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("ActualTPS: %f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func init() {
}

type Game struct {
	w       World
	player1 *ebiten.Image
	player2 *ebiten.Image
	ball1   *ebiten.Image
	ball2   *ebiten.Image
	health  *ebiten.Image
}

func loadImage(str string) *ebiten.Image {
	file, err := os.Open(str)
	defer file.Close()
	Check(err)

	img, _, err := image.Decode(file)
	Check(err)
	return ebiten.NewImageFromImage(img)
}

func main() {
	var g Game
	g.ball1 = loadImage("sprites/ball1.png")
	g.ball2 = loadImage("sprites/ball2.png")
	g.player1 = loadImage("sprites/player1.png")
	g.player2 = loadImage("sprites/player2.png")
	g.health = loadImage("sprites/health.png")
	ebiten.SetWindowSize(460, 460)
	ebiten.SetWindowTitle("Viewer")
	ebiten.SetWindowPosition(10, 1080-470)
	err := ebiten.RunGame(&g)
	Check(err)
}
