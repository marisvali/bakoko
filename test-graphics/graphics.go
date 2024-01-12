package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
	"math"
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
	// Get keyboard input.
	var pressedKeys []ebiten.Key
	pressedKeys = inpututil.AppendPressedKeys(pressedKeys)
	// Choose which is the active player based on Alt being pressed.
	//playerInput := PlayerInput{}
	//playerInput.MoveLeft = slices.Contains(pressedKeys, ebiten.KeyA)
	//playerInput.MoveUp = slices.Contains(pressedKeys, ebiten.KeyW)
	//playerInput.MoveDown = slices.Contains(pressedKeys, ebiten.KeyS)
	//playerInput.MoveRight = slices.Contains(pressedKeys, ebiten.KeyD)
	step := U(1)
	if slices.Contains(pressedKeys, ebiten.KeyA) {
		g.c.Center.X.Subtract(step)
	}
	if slices.Contains(pressedKeys, ebiten.KeyD) {
		g.c.Center.X.Add(step)
	}
	if slices.Contains(pressedKeys, ebiten.KeyW) {
		g.c.Center.Y.Subtract(step)
	}
	if slices.Contains(pressedKeys, ebiten.KeyS) {
		g.c.Center.Y.Add(step)
	}

	var justPressedKeys []ebiten.Key
	justPressedKeys = inpututil.AppendJustPressedKeys(justPressedKeys)

	// Get mouse input.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		//x, y := ebiten.CursorPosition()
	}

	//g.w.Step(&input)
	//var w World
	//g.peer.getWorld(&w)
	//input.SerializeToFile("input.bin")
	//TouchFile("input-ready")
	//WaitForFile("world-ready")
	//g.w.DeserializeFromFile("world.bin")
	return nil
}

func DrawSprite(screen *ebiten.Image, img *ebiten.Image, pos Pt) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X.ToFloat64(), pos.Y.ToFloat64())
	screen.DrawImage(img, op)
}

func DrawPixel(screen *ebiten.Image, x, y int, color color.Color) {
	size := 0
	for ax := x - size; ax <= x+size; ax++ {
		for ay := y - size; ay <= y+size; ay++ {
			screen.Set(ax, ay, color)
		}
	}
}

func DrawPixel2(screen *ebiten.Image, pt Pt, color color.Color) {
	x := int(WorldToScreen(pt.X))
	y := int(WorldToScreen(pt.Y))
	size := 0
	for ax := x - size; ax <= x+size; ax++ {
		for ay := y - size; ay <= y+size; ay++ {
			screen.Set(ax, ay, color)
		}
	}
}
func DrawLine2(screen *ebiten.Image, x1, y1, x2, y2 float64, color color.Color) {
	if math.Abs(x1-x2) > math.Abs(y1-y2) {
		startX := int(math.Min(x1, x2))
		endX := int(math.Round(math.Max(x1, x2)))
		for x := startX; x <= endX; x++ {
			factor := float64(x) / float64(endX-startX)
			y := int(y1 + factor*(y2-y1))
			//screen.Set(x, y, colorPrimary)
			DrawPixel(screen, x, y, color)
		}
	} else {
		startY := int(math.Min(y1, y2))
		endY := int(math.Round(math.Max(y1, y2)))
		for y := startY; y <= endY; y++ {
			factor := float64(y) / float64(endY-startY)
			x := int(x1 + factor*(x2-x1))
			//screen.Set(x, y, colorPrimary)
			DrawPixel(screen, x, y, color)
		}
	}
}

func DrawLine(screen *ebiten.Image, l Line, color color.Color) {
	x1 := WorldToScreen(l.Start.X)
	y1 := WorldToScreen(l.Start.Y)
	x2 := WorldToScreen(l.End.X)
	y2 := WorldToScreen(l.End.Y)
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	/*
		int dx = x2 - x1;
		int dy = y2 - y1;

		for (int x = x1; x < x2; x++)
		{
			int y = y1 + dy * (x - x1) / dx;
			//[y*canvas.Width+x] converts the 2d array index to a 1d array index
			canvasData[y * canvas.Width + x] = Color.Black;
		}
	*/
	if math.Abs(x1-x2) > math.Abs(y1-y2) {
		startX := int(math.Min(x1, x2))
		endX := int(math.Round(math.Max(x1, x2)))
		for x := startX; x <= endX; x++ {
			factor := float64(x-startX) / float64(endX-startX)
			y := int(y1 + factor*(y2-y1))
			//screen.Set(x, y, colorPrimary)
			DrawPixel(screen, x, y, color)
		}
	} else {
		startY := int(math.Min(y1, y2))
		endY := int(math.Round(math.Max(y1, y2)))
		for y := startY; y <= endY; y++ {
			factor := float64(y-startY) / float64(endY-startY)
			x := int(x1 + factor*(x2-x1))
			//screen.Set(x, y, colorPrimary)
			DrawPixel(screen, x, y, color)
		}
	}
}

func DrawCircle2(screen *ebiten.Image, x, y float64, r float64, color color.Color) {
	// calculates the minimun angle between two pixels in a diagonal.
	// you can multiply minAngle by a security factor like 0.9 just to be sure you wont have empty pixels in the circle
	minAngle := math.Acos(1.0 - 1.0/r)

	for angle := float64(0); angle <= 360.0; angle += minAngle {
		x1 := r * math.Cos(angle)
		y1 := r * math.Sin(angle)
		DrawPixel(screen, int(x+x1), int(y+y1), color)
	}
}

func WorldToScreen(val Int) float64 {
	return val.ToFloat64() / Unit * 10
}

func DrawCircle(screen *ebiten.Image, c Circle, color color.Color) {
	x := WorldToScreen(c.Center.X)
	y := WorldToScreen(c.Center.Y)
	r := WorldToScreen(c.Diameter) / 2
	// calculates the minimun angle between two pixels in a diagonal.
	// you can multiply minAngle by a security factor like 0.9 just to be sure you wont have empty pixels in the circle
	minAngle := math.Acos(1.0 - 1.0/r)

	for angle := float64(0); angle <= 360.0; angle += minAngle {
		x1 := r * math.Cos(angle)
		y1 := r * math.Sin(angle)
		DrawPixel(screen, int(x+x1), int(y+y1), color)
	}
}

func DrawSquare(screen *ebiten.Image, s Square, color color.Color) {
	halfSize := s.Size.DivBy(I(2)).Plus(s.Size.Mod(I(2)))

	// square corners
	upperLeftCorner := Pt{s.Center.X.Minus(halfSize), s.Center.Y.Minus(halfSize)}
	lowerLeftCorner := Pt{s.Center.X.Minus(halfSize), s.Center.Y.Plus(halfSize)}
	upperRightCorner := Pt{s.Center.X.Plus(halfSize), s.Center.Y.Minus(halfSize)}
	lowerRightCorner := Pt{s.Center.X.Plus(halfSize), s.Center.Y.Plus(halfSize)}

	DrawLine(screen, Line{upperLeftCorner, upperRightCorner}, color)
	DrawLine(screen, Line{upperLeftCorner, lowerLeftCorner}, color)
	DrawLine(screen, Line{lowerLeftCorner, lowerRightCorner}, color)
	DrawLine(screen, Line{lowerRightCorner, upperRightCorner}, color)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(colorNeutralLight1)
	//line := Line{IPt(200, 100), IPt(100, 300)}
	//DrawLine(screen, line, colorSecondary)
	//circle := Circle{IPt(300, 200), I(30)}
	//DrawCircle(screen, circle, colorPrimary)

	//l1 := Line{IPt(30, 100), IPt(400, 300)}
	//l2 := Line{IPt(250, 100), IPt(250, 300)}
	//DrawLine(screen, l1, colorPrimary)
	//DrawLine(screen, l2, colorPrimaryDark1)
	//
	//intersects, pt := LineVerticalLineIntersection(l1, l2)
	////intersects, pt := LineHorizontalLineIntersection(l1, l2)
	//if intersects {
	//	DrawPixel2(screen, pt, colorSecondary)
	//}

	// Debug line-circle intersection.
	//l := Line{IPt(130, 60), IPt(500, 500)}
	//c := Circle{IPt(200, 200), I(150)}
	//DrawLine(screen, l, colorPrimary)
	//DrawCircle(screen, c, colorPrimaryDark1)
	//DrawPixel2(screen, c.Center, colorSecondaryLight2)
	//
	//intersects, pt := LineCircleIntersection(l, c)
	//if intersects {
	//	DrawPixel2(screen, pt, colorSecondary)
	//}

	// Debug circle-square intersection.
	c := g.c
	s := g.s

	DrawCircle(screen, c, colorPrimaryDark1)
	DrawSquare(screen, s, colorPrimaryDark1)
	DrawPixel2(screen, s.Center, colorSecondary)
	intersects, circlePositionAtCollision, _, debugInfo := CircleSquareCollision(c.Center, s.Center, c.Diameter, s)

	if intersects {
		c.Center = circlePositionAtCollision
		DrawCircle(screen, c, colorPrimaryLight2)
		DrawPixel2(screen, circlePositionAtCollision, colorSecondaryDark2)
	}

	for _, l := range debugInfo.Lines {
		DrawLine(screen, l, color.RGBA{0, 0, 255, 255})
	}
	for _, c := range debugInfo.Circles {
		DrawCircle(screen, c, color.RGBA{0, 0, 255, 255})
	}

	for _, p := range debugInfo.Points {
		DrawPixel2(screen, p, color.RGBA{255, 255, 0, 255})
	}

	//img1 := ebiten.NewImage(50, 50)
	//img1.Fill(colorPrimary)
	//op := &ebiten.DrawImageOptions{}
	//op.GeoM.Translate(Real(g.w.Player1.Center.X), Real(g.w.Player1.Center.Y))
	//screen.DrawImage(img1, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("x: %d y: %d", g.c.Center.X.DivBy(U(1)).ToInt64(), g.c.Center.Y.DivBy(U(1)).ToInt64()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

type Game struct {
	c Circle
	s Square
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
	//g.c = Circle{IPt(70, 200), I(100)}
	//g.c = Circle{UPt(270, 30), U(100)}
	//g.s = Square{UPt(300, 200), U(100)}
	//g.c = Circle{UPt(270, 30), U(10)}
	//g.s = Square{UPt(300, 200), U(10)}
	g.c = Circle{UPt(50, 50), U(10)}
	g.s = Square{UPt(50, 50), U(10)}
	//ebiten.SetWindowSize(460, 460)
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("Viewer")
	ebiten.SetWindowPosition(10, 1080-470)
	err := ebiten.RunGame(&g)
	Check(err)
}
