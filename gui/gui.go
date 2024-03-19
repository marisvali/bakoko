package gui

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"math"
	"os"
	. "playful-patterns.com/bakoko/ints"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
	"slices"
	"sync"
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

func playerWasJustHit(player Player, prevHealth Int) bool {
	return player.Health.Lt(prevHealth)
}

func (g *Gui) UpdateGameOngoing() {
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
	playerInput.Reload = slices.Contains(justPressedKeys, ebiten.KeyR)

	if slices.Contains(justPressedKeys, ebiten.KeyEscape) {
		g.state = GamePaused
	}

	if g.folderWatcher.FolderContentsChanged() {
		g.loadGuiData()
	}

	// Get mouse input.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		playerInput.Shoot = true
		x, y := ebiten.CursorPosition()
		// Translate from screen coordinates to in-world-main units.
		playerInput.ShootPt.X = g.ScreenToWorld(x)
		playerInput.ShootPt.Y = g.ScreenToWorld(y)
	}

	if g.SyncWithWorld(playerInput) {
		// React to updates.
		// Player1 was just hit.
		if playerWasJustHit(g.w.Player1, g.player1PreviousHealth) {
			g.hitAnimation1 = 255
		}
		g.player1PreviousHealth = g.w.Player1.Health

		// Player2 was just hit.
		if playerWasJustHit(g.w.Player2, g.player2PreviousHealth) {
			g.hitAnimation2 = 255
		}
		g.player2PreviousHealth = g.w.Player2.Health

		if g.w.Player1.Health.Eq(ZERO) {
			g.state = GameLost
			g.gameOverAnimation = -500
		}

		if g.w.Player2.Health.Eq(ZERO) {
			g.state = GameWon
			g.gameOverAnimation = -500
		}
	}

	if g.hitAnimation1 > 0 {
		g.hitAnimation1 -= 10
	}

	if g.hitAnimation2 > 0 {
		g.hitAnimation2 -= 10
	}
}

func (g *Gui) UpdateGamePaused() {
	// Get keyboard input.
	var pressedKeys []ebiten.Key
	pressedKeys = inpututil.AppendPressedKeys(pressedKeys)

	playerInput := PlayerInput{}
	playerInput.Pause = true

	var justPressedKeys []ebiten.Key
	justPressedKeys = inpututil.AppendJustPressedKeys(justPressedKeys)
	if slices.Contains(justPressedKeys, ebiten.KeyR) {
		playerInput.Reload = true
		g.state = GameOngoing
	}

	unpause := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	unpauseKeys := []ebiten.Key{ebiten.KeyA, ebiten.KeyW, ebiten.KeyS, ebiten.KeyD}
	for i := range unpauseKeys {
		if slices.Contains(justPressedKeys, unpauseKeys[i]) {
			unpause = true
			break
		}
	}

	if unpause {
		playerInput.Pause = false
		g.state = GameOngoing
	}

	if g.folderWatcher.FolderContentsChanged() {
		g.loadGuiData()
	}

	g.SyncWithWorld(playerInput)
}

func (g *Gui) UpdateGameWon() {
	// Get keyboard input.
	var pressedKeys []ebiten.Key
	pressedKeys = inpututil.AppendPressedKeys(pressedKeys)

	playerInput := PlayerInput{}
	//playerInput.Pause = true

	var justPressedKeys []ebiten.Key
	justPressedKeys = inpututil.AppendJustPressedKeys(justPressedKeys)
	if slices.Contains(justPressedKeys, ebiten.KeyR) {
		playerInput.Reload = true
		g.state = GameOngoing
	}

	if g.folderWatcher.FolderContentsChanged() {
		g.loadGuiData()
	}

	g.SyncWithWorld(playerInput)

	if g.hitAnimation1 > 0 {
		g.hitAnimation1 -= 10
	}
	if g.hitAnimation2 > 0 {
		g.hitAnimation2 -= 10
	}
}

func (g *Gui) UpdatePlayback() {
	// Get keyboard input.
	var justPressedKeys []ebiten.Key
	justPressedKeys = inpututil.AppendJustPressedKeys(justPressedKeys)

	//if slices.Contains(justPressedKeys, ebiten.KeyEscape) {
	//	g.state = GamePaused
	//}

	// Get mouse input.
	g.leftButtonClicked = inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		g.leftButtonPressed = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		g.leftButtonPressed = false
	}
	g.mousePosX, g.mousePosY = ebiten.CursorPosition()

	if g.folderWatcher.FolderContentsChanged() {
		g.loadGuiData()
	}

	if g.targetFrame >= 0 {
		// Rewind.
		var initialPlayerInput PlayerInput
		initialPlayerInput.Reload = true
		initialPlayerInput.Pause = true
		g.SyncWithWorld(initialPlayerInput)
		for i := 0; i < g.targetFrame; i++ {
			g.SyncWithWorld(g.playerInputs[i])
		}
		g.frameIdx = g.targetFrame

		g.player1PreviousHealth = g.w.Player1.Health
		g.player2PreviousHealth = g.w.Player2.Health
		g.targetFrame = -1
	}

	var playerInput PlayerInput
	if g.frameIdx < len(g.playerInputs) {
		playerInput = g.playerInputs[g.frameIdx]
	}
	g.frameIdx++

	if g.SyncWithWorld(playerInput) {
		// React to updates.
		// Player1 was just hit.
		if playerWasJustHit(g.w.Player1, g.player1PreviousHealth) {
			g.hitAnimation1 = 255
		}
		g.player1PreviousHealth = g.w.Player1.Health

		// Player2 was just hit.
		if playerWasJustHit(g.w.Player2, g.player2PreviousHealth) {
			g.hitAnimation2 = 255
		}
		g.player2PreviousHealth = g.w.Player2.Health
	}

	if g.hitAnimation1 > 0 {
		g.hitAnimation1 -= 10
	}

	if g.hitAnimation2 > 0 {
		g.hitAnimation2 -= 10
	}

	// Update the AI if there is one.
	if g.aiRunner != nil {
		g.aiRunner.Step()
	}

	// Update the world if there is one.
	if g.worldRunner != nil {
		g.worldRunner.Step()
	}
}

func (g *Gui) UpdateGameLost() {
	// Get keyboard input.
	var pressedKeys []ebiten.Key
	pressedKeys = inpututil.AppendPressedKeys(pressedKeys)

	playerInput := PlayerInput{}
	//playerInput.Pause = true

	var justPressedKeys []ebiten.Key
	justPressedKeys = inpututil.AppendJustPressedKeys(justPressedKeys)
	if slices.Contains(justPressedKeys, ebiten.KeyR) {
		playerInput.Reload = true
		g.state = GameOngoing
	}

	if g.folderWatcher.FolderContentsChanged() {
		g.loadGuiData()
	}

	g.SyncWithWorld(playerInput)

	if g.hitAnimation1 > 0 {
		g.hitAnimation1 -= 10
	}
	if g.hitAnimation2 > 0 {
		g.hitAnimation2 -= 10
	}
}

func (g *Gui) Update() error {
	if g.state == GameOngoing {
		g.UpdateGameOngoing()
	} else if g.state == GamePaused {
		g.UpdateGamePaused()
	} else if g.state == GameWon {
		g.UpdateGameWon()
	} else if g.state == GameLost {
		g.UpdateGameLost()
	} else if g.state == Playback {
		g.UpdatePlayback()
	}
	return nil
}

func (g *Gui) SyncWithWorld(input PlayerInput) bool {
	// Here I want to block but only if there's a connection.
	// If a connection cannot be established, or the send or get fails, or
	// there is a timeout, I want to go ahead.
	// I will try to get the connection back at every update, but I don't want
	// to permanently block my GUI if a connection cannot be established.
	if err := g.worldProxy.Connect(); err != nil {
		return false // Nevermind, try again next frame.
	}

	var err error
	if g.w, err = g.worldProxy.GetWorld(); err != nil {
		return false // Nevermind, try again next frame.
	}

	if err = g.worldProxy.SendInput(&input); err != nil {
		return false // Nevermind, try again next frame.
	}

	return true
}

//func DrawSprite(screen *ebiten.Image, img *ebiten.Image, pos Pt) {
//	op := &ebiten.DrawImageOptions{}
//	op.GeoM.Translate(pos.X.ToFloat64(), pos.Y.ToFloat64())
//	screen.DrawImage(img, op)
//}

func (g *Gui) WorldToScreen(val Int) float64 {
	return val.ToFloat64() / Unit * g.data.ScaleFactor
}

func (g *Gui) ScreenToWorld(val int) Int {
	return U(int(float64(val) / g.data.ScaleFactor))
}

func (g *Gui) DrawSprite(img *ebiten.Image,
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

func (g *Gui) DrawSprite2(img *ebiten.Image,
	x float64, y float64, targetWidth float64, targetHeight float64) {
	op := &ebiten.DrawImageOptions{}

	// Resize image to fit the target size we want to draw.
	// This kind of scaling is very useful during development when the final
	// sizes are not decided, and thus it's impossible to have final sprites.
	// For an actual release, scaling should be avoided.
	imgSize := img.Bounds().Size()
	newDx := targetWidth / float64(imgSize.X)
	newDy := targetHeight / float64(imgSize.Y)
	op.GeoM.Scale(newDx, newDy)

	op.GeoM.Translate(x, y)

	g.screen.DrawImage(img, op)
}

func (g *Gui) DrawSprite3(img *ebiten.Image,
	x float64, y float64, targetWidth float64, targetHeight float64, alpha float32) {
	op := &ebiten.DrawImageOptions{}

	// Resize image to fit the target size we want to draw.
	// This kind of scaling is very useful during development when the final
	// sizes are not decided, and thus it's impossible to have final sprites.
	// For an actual release, scaling should be avoided.
	imgSize := img.Bounds().Size()
	newDx := targetWidth / float64(imgSize.X)
	newDy := targetHeight / float64(imgSize.Y)
	op.GeoM.Scale(newDx, newDy)

	op.GeoM.Translate(x, y)
	op.ColorScale.SetR(alpha)
	op.ColorScale.SetG(alpha)
	op.ColorScale.SetB(alpha)
	op.ColorScale.SetA(alpha)

	g.screen.DrawImage(img, op)
}

func (g *Gui) DrawPlayer(
	playerImage *ebiten.Image,
	ballImage *ebiten.Image,
	healthImage *ebiten.Image,
	player *Player) {

	// Draw the player sprite.
	x := g.WorldToScreen(player.Bounds.Center.X) + g.data.PlayerOffsetX*g.data.ScaleFactor
	y := g.WorldToScreen(player.Bounds.Center.Y) + g.data.PlayerOffsetY*g.data.ScaleFactor
	diam := g.WorldToScreen(player.Bounds.Diameter) + g.data.PlayerOffsetSize*g.data.ScaleFactor
	g.DrawSprite(playerImage, x, y, diam)

	// Draw a small sprite for each ball that the player has.
	realDiam := g.WorldToScreen(player.Bounds.Diameter)
	for idx := int64(0); idx < player.NBalls.ToInt64(); idx++ {
		col := float64(idx / 5)
		row := float64(idx % 5)
		smallX := x - realDiam/2 - 10*g.data.ScaleFactor - col*10*g.data.ScaleFactor
		smallY := y + (float64(row*12)+5)*g.data.ScaleFactor - realDiam/2
		smallDiam := float64(10) * g.data.ScaleFactor
		g.DrawSprite(ballImage, smallX, smallY, smallDiam)
	}

	// Draw a small sprite for each health point that the player has.
	startX := x - realDiam/2
	fullWidth := (float64(6 * 12)) * g.data.ScaleFactor
	actualWidth := (float64(player.Health.ToInt64() * 12)) * g.data.ScaleFactor
	startX += (fullWidth - actualWidth) / 2
	for idx := int64(0); idx < player.Health.ToInt64(); idx++ {
		smallX := startX + (float64(idx*12))*g.data.ScaleFactor
		smallY := y - realDiam/2 - 10*g.data.ScaleFactor
		smallDiam := float64(10) * g.data.ScaleFactor
		g.DrawSprite(healthImage, smallX, smallY, smallDiam)
	}

	// Draw actual bounds, for debugging purposes.
	if g.data.DrawDebugGraphics {
		g.DrawCircle(player.Bounds, color.White)
	}
}

func (g *Gui) DrawCircle(c Circle, color color.Color) {
	x := g.WorldToScreen(c.Center.X)
	y := g.WorldToScreen(c.Center.Y)
	r := g.WorldToScreen(c.Diameter) / 2
	// calculates the minimun angle between two pixels in a diagonal.
	// you can multiply minAngle by a security factor like 0.9 just to be sure you wont have empty pixels in the circle
	minAngle := math.Acos(1.0 - 1.0/r)

	for angle := float64(0); angle <= 360.0; angle += minAngle {
		x1 := r * math.Cos(angle)
		y1 := r * math.Sin(angle)
		DrawPixel(g.screen, int(x+x1), int(y+y1), color)
	}
}

func DrawPixel(screen *ebiten.Image, x, y int, color color.Color) {
	size := 0
	for ax := x - size; ax <= x+size; ax++ {
		for ay := y - size; ay <= y+size; ay++ {
			screen.Set(ax, ay, color)
		}
	}
}

func (g *Gui) DrawFilledSquare(screen *ebiten.Image, s Square, col color.Color) {
	size := int(g.WorldToScreen(s.Size))
	x1 := int(g.WorldToScreen(s.Center.X)) - size/2
	y1 := int(g.WorldToScreen(s.Center.Y)) - size/2
	x2 := x1 + size
	y2 := y1 + size
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			screen.Set(x, y, col)
		}
	}

	//if g.filledSquare == nil {
	//	g.filledSquare = ebiten.NewImage(int(size), int(size))
	//}
	//g.filledSquare.Fill(col)
	//op := &ebiten.DrawImageOptions{}
	//op.GeoM.Translate(x, y)
	//screen.DrawImage(g.filledSquare, op)
}

func (g *Gui) Draw(screen *ebiten.Image) {
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

	// Player1
	if g.w.Player1.State.Eq(PlayerStunned) {
		g.DrawPlayer(g.player1Hit, g.ball1, g.health, &g.w.Player1)
	} else {
		g.DrawPlayer(g.player1, g.ball1, g.health, &g.w.Player1)
	}

	// Player2
	if g.w.Player2.State.Eq(PlayerStunned) {
		g.DrawPlayer(g.player2Hit, g.ball2, g.health, &g.w.Player2)
	} else {
		g.DrawPlayer(g.player2, g.ball2, g.health, &g.w.Player2)
	}

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

	// Draw instructional text.
	var textHeight float64 = 50
	g.DrawSprite2(g.textBackground, 0,
		float64(screen.Bounds().Dy())-textHeight-float64(g.data.PlaybackBarHeight),
		float64(screen.Bounds().Dx()),
		textHeight)
	var message string
	if g.state == GameOngoing {
		message = "Defeat your opponent! Press WASD to move, left click to shoot, R to restart, ESC to pause, move or shoot to unpause."
	} else if g.state == GamePaused {
		message = "Defeat your opponent! Press WASD to move, left click to shoot, R to restart, ESC to pause, move or shoot to unpause."
	} else if g.state == GameWon {
		message = "You won, congratulations! Press R to play again."
	} else if g.state == GameLost {
		message = "You lost. Press R to play again."
	} else if g.state == Playback {
		message = fmt.Sprintf("Playing back frame %d / %d", g.frameIdx, len(g.playerInputs))
	} else {
		Check(fmt.Errorf("unhandled game state: %d", g.state))
	}

	textSize := text.BoundString(g.defaultFont, message)
	textX := screen.Bounds().Min.X + (screen.Bounds().Dx()-textSize.Dx())/2
	textY := screen.Bounds().Max.Y - (int(textHeight)-textSize.Dy())/2 - g.data.PlaybackBarHeight
	text.Draw(screen, message, g.defaultFont, textX, textY, colorHex(0x000000))

	if g.state == GamePaused {
		message = "PAUSED"
		xMargin := 60
		textSize := text.BoundString(g.defaultFont, message)

		textX1 := xMargin
		textY := screen.Bounds().Max.Y - (int(textHeight)-textSize.Dy())/2 - g.data.PlaybackBarHeight
		text.Draw(screen, message, g.defaultFont, textX1, textY, colorHex(0xee005a))

		textX2 := screen.Bounds().Min.X + (screen.Bounds().Dx() - textSize.Dx() - xMargin)
		text.Draw(screen, message, g.defaultFont, textX2, textY, colorHex(0xee005a))
	}

	// Hit animations.
	if g.hitAnimation1 > 0 {
		//col := color.RGBA{255, 0, 0, uint8(g.hitAnimation)}
		//g.DrawFilledSquare(screen, Square{Pt{U(500), U(500)}, U(50)}, col)

		op := &ebiten.DrawImageOptions{}
		// Scale image to cover the entire screen.
		imgSize := g.hit.Bounds().Size()
		targetSize := screen.Bounds().Size()
		targetSize.Y -= int(textHeight) - g.data.PlaybackBarHeight
		newDx := float64(targetSize.X) / float64(imgSize.X)
		newDy := float64(targetSize.Y) / float64(imgSize.Y)
		op.GeoM.Scale(newDx, newDy)
		op.GeoM.Translate(0, 0)
		colorScale := float32(g.hitAnimation1) / float32(255)
		op.ColorScale.SetR(colorScale)
		op.ColorScale.SetG(colorScale)
		op.ColorScale.SetB(colorScale)
		op.ColorScale.SetA(colorScale)
		screen.DrawImage(g.hit, op)
	}
	if g.hitAnimation2 > 0 {
		//col := color.RGBA{255, 0, 0, uint8(g.hitAnimation)}
		//g.DrawFilledSquare(screen, Square{Pt{U(500), U(500)}, U(50)}, col)

		op := &ebiten.DrawImageOptions{}
		// Scale image to cover the entire screen.
		imgSize := g.hitGood.Bounds().Size()
		targetSize := screen.Bounds().Size()
		targetSize.Y -= int(textHeight) - g.data.PlaybackBarHeight
		newDx := float64(targetSize.X) / float64(imgSize.X)
		newDy := float64(targetSize.Y) / float64(imgSize.Y)
		op.GeoM.Scale(newDx, newDy)
		op.GeoM.Translate(0, 0)
		colorScale := float32(g.hitAnimation2) / float32(255)
		op.ColorScale.SetR(colorScale)
		op.ColorScale.SetG(colorScale)
		op.ColorScale.SetB(colorScale)
		op.ColorScale.SetA(colorScale)
		screen.DrawImage(g.hitGood, op)
	}

	if g.gameOverAnimation > 0 {
		alpha := float32(g.gameOverAnimation) / 255
		if g.state == GameWon {
			targetSizeX := float64(screen.Bounds().Size().X)
			targetSizeY := float64(screen.Bounds().Size().Y) - textHeight - float64(g.data.PlaybackBarHeight)
			g.DrawSprite3(g.won, 0, 0, targetSizeX, targetSizeY, alpha)
		} else if g.state == GameLost {
			targetSizeX := float64(screen.Bounds().Size().X)
			targetSizeY := float64(screen.Bounds().Size().Y) - textHeight - float64(g.data.PlaybackBarHeight)
			g.DrawSprite3(g.lost, 0, 0, targetSizeX, targetSizeY, alpha)
		}
	}

	if g.gameOverAnimation > -1000 && g.gameOverAnimation < 255 {
		g.gameOverAnimation += 10
	}

	if g.state == Playback {
		g.DrawPlaybackBar(screen)
	}

	// Debug geometry.
	for i := range g.debugInfo {
		info := g.GetDebugInfo(i)
		for _, p := range info.Points {
			g.DrawFilledSquare(screen, Square{p.Pos, p.Size}, p.Col)
		}
	}

	//img1 := ebiten.NewImage(50, 50)
	//img1.Fill(colorPrimary)
	//op := &ebiten.DrawImageOptions{}
	//op.GeoM.Translate(Real(g.w.Player1.Pos.X), Real(g.w.Player1.Pos.Y))
	//screen.DrawImage(img1, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("ActualTPS: %f", ebiten.ActualTPS()))

	g.screen = nil
}

func (g *Gui) DrawPlaybackBar(screen *ebiten.Image) {
	g.DrawSprite2(
		g.textBackground, 0,
		float64(screen.Bounds().Dy()-g.data.PlaybackBarHeight),
		float64(screen.Bounds().Dx()),
		float64(g.data.PlaybackBarHeight))

	if g.leftButtonPressed {
		g.DrawSprite(
			g.ball1,
			float64(g.mousePosX),
			float64(g.mousePosY),
			30)
	}

	var x, y, width, height float64
	x = 10
	y = float64(screen.Bounds().Dy()-g.data.PlaybackBarHeight) + 10
	width = float64(screen.Bounds().Dx()) - 300
	height = float64(g.data.PlaybackBarHeight) - 20
	mx, my := float64(g.mousePosX), float64(g.mousePosY)

	g.DrawSprite2(g.health, x, y, width, height)
	if g.leftButtonClicked &&
		mx >= x && mx <= (x+width) &&
		my >= y && my <= (y+height) {
		g.DrawSprite(
			g.ball2,
			float64(g.mousePosX),
			float64(g.mousePosY),
			30)

		factor := (mx - x) / width
		g.targetFrame = int(factor * float64(len(g.playerInputs)))
	}
}

func (g *Gui) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

type GuiData struct {
	ScaleFactor       float64
	WindowWidth       int
	WindowHeight      int
	PlayerOffsetX     float64
	PlayerOffsetY     float64
	PlayerOffsetSize  float64
	DrawDebugGraphics bool
	PlaybackBarHeight int
}

type Gui struct {
	w              *World
	worldProxy     WorldProxy
	painters       []PainterProxy
	player1        *ebiten.Image
	player1Hit     *ebiten.Image
	player2        *ebiten.Image
	player2Hit     *ebiten.Image
	ball1          *ebiten.Image
	ball2          *ebiten.Image
	health         *ebiten.Image
	obstacle       *ebiten.Image
	background     *ebiten.Image
	screen         *ebiten.Image
	hit            *ebiten.Image
	hitGood        *ebiten.Image
	textBackground *ebiten.Image
	won            *ebiten.Image
	lost           *ebiten.Image
	data           GuiData
	times          []time.Time
	filledSquare   *ebiten.Image
	debugInfo      []DebugInfo
	debugInfoMutex []sync.Mutex
	folderWatcher  FolderWatcher
	hitAnimation1  int
	hitAnimation2  int
	// The UI is responsible for keeping track of state changes that are
	// relevant for it.
	player1PreviousHealth Int
	player2PreviousHealth Int
	state                 GameState
	defaultFont           font.Face
	gameOverAnimation     int
	playerInputs          []PlayerInput
	frameIdx              int
	leftButtonClicked     bool
	leftButtonPressed     bool
	mousePosX             int
	mousePosY             int
	targetFrame           int
	worldRunner           *WorldRunner
	aiRunner              *AiRunner
}

type GameState int64

const (
	GameOngoing GameState = iota
	GamePaused
	GameWon
	GameLost
	Playback
)

func loadImage(str string) *ebiten.Image {
	file, err := os.Open(str)
	defer file.Close()
	Check(err)

	img, _, err := image.Decode(file)
	Check(err)
	if err != nil {
		return nil
	}

	return ebiten.NewImageFromImage(img)
}

//func loadJSON(filename string, v any) {
//	file, err := os.Open(filename)
//	Check(err)
//	bytes, err := io.ReadAll(file)
//	Check(err)
//	err = json.Unmarshal(bytes, v)
//	Check(err)
//}
//
//func (g *Gui) gameDataChangedOnDisk() bool {
//	files, err := os.ReadDir("gui-data")
//	Check(err)
//	if len(files) != len(g.times) {
//		g.times = make([]time.Time, len(files))
//	}
//	changed := false
//	for idx, file := range files {
//		info, err := file.Info()
//		Check(err)
//		if g.times[idx] != info.ModTime() {
//			changed = true
//			g.times[idx] = info.ModTime()
//		}
//	}
//	return changed
//}

func (g *Gui) loadGuiData() {
	// Read from the disk over and over until a full read is possible.
	// This repetition is meant to avoid crashes due to reading files
	// while they are still being written.
	// It's a hack but possibly a quick and very useful one.
	CheckCrashes = false
	for {
		CheckFailed = nil
		g.ball1 = loadImage("gui-data/ball1.png")
		g.ball2 = loadImage("gui-data/ball2.png")
		g.player1 = loadImage("gui-data/player1.png")
		g.player1Hit = loadImage("gui-data/player1-hit.png")
		g.player2 = loadImage("gui-data/player2.png")
		g.player2Hit = loadImage("gui-data/player2-hit.png")
		g.health = loadImage("gui-data/health.png")
		g.obstacle = loadImage("gui-data/obstacle.png")
		g.background = loadImage("gui-data/background.png")
		g.hit = loadImage("gui-data/hit.png")
		g.hitGood = loadImage("gui-data/hit-good.png")
		g.textBackground = loadImage("gui-data/text-background.png")
		g.won = loadImage("gui-data/won.png")
		g.lost = loadImage("gui-data/lost.png")
		if g.state == Playback {
			LoadJSON("gui-data/gui-playback.json", &g.data)
		} else {
			LoadJSON("gui-data/gui.json", &g.data)
		}

		if CheckFailed == nil {
			break
		}
	}
	CheckCrashes = true

	ebiten.SetWindowSize(g.data.WindowWidth, g.data.WindowHeight)
	ebiten.SetWindowTitle("Viewer")
	ebiten.SetWindowPosition(10, 1080-10-g.data.WindowHeight)
}

func (g *Gui) SetDebugInfo(i int, info DebugInfo) {
	g.debugInfoMutex[i].Lock()
	g.debugInfo[i] = info.Clone() // Must to deep copy here.
	g.debugInfoMutex[i].Unlock()
}

func (g *Gui) GetDebugInfo(i int) DebugInfo {
	g.debugInfoMutex[i].Lock()
	info := g.debugInfo[i].Clone() // Must to deep copy here.
	g.debugInfoMutex[i].Unlock()
	return info
}

func (g *Gui) UpdateDebugInfo(i int) {
	for {
		info := g.painters[i].GetPaintData() // Block.
		g.SetDebugInfo(i, info)
	}
}

func (g *Gui) AddPainter(endpoint string) {
	var p PainterProxyTcpIp
	p.Endpoint = endpoint
	g.painters = append(g.painters, &p)
	g.debugInfo = append(g.debugInfo, DebugInfo{})
	g.debugInfoMutex = append(g.debugInfoMutex, sync.Mutex{})
	i := len(g.painters) - 1
	go g.UpdateDebugInfo(i)
}

func (g *Gui) Init(worldProxy WorldProxy, worldRunner *WorldRunner, aiRunner *AiRunner, recordingFile string) {
	g.worldProxy = worldProxy
	g.worldRunner = worldRunner
	g.aiRunner = aiRunner

	g.frameIdx = 0
	g.targetFrame = -1
	if recordingFile == "" {
		g.state = GamePaused
	} else {
		g.state = Playback
		g.playerInputs = deserializeInputs(recordingFile)
	}

	g.folderWatcher.Folder = "gui-data"
	//g.AddPainter(os.Args[2])
	//g.AddPainter(os.Args[3])
	g.loadGuiData()

	// Load the Arial font
	fontData, err := opentype.Parse(goregular.TTF)
	Check(err)

	g.defaultFont, err = opentype.NewFace(fontData, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	Check(err)
}

func deserializeInputs(filename string) []PlayerInput {
	var inputs []PlayerInput
	buf := bytes.NewBuffer(ReadFile(filename))
	var lenInputs Int
	Deserialize(buf, &lenInputs)
	inputs = make([]PlayerInput, lenInputs.ToInt64())
	Deserialize(buf, inputs)
	return inputs
}

func RunGui(worldProxy WorldProxy) {
	var g Gui
	g.Init(worldProxy, nil, nil, "")

	// Start the game.
	err := ebiten.RunGame(&g)
	Check(err)
}

func RunGuiPlayback(recordingFile string) {
	var worldAiProxy WorldPlayerProxy  // Connects the world and AI.
	var worldGuiProxy WorldPlayerProxy // Connects the world and GUI.
	var worldRunner WorldRunner
	var aiRunner AiRunner
	aiRunner.Initialize(&worldAiProxy)
	worldRunner.Initialize(&worldGuiProxy, &worldAiProxy)

	var g Gui
	g.Init(&worldGuiProxy, &worldRunner, &aiRunner, recordingFile)

	// Start the game.
	err := ebiten.RunGame(&g)
	Check(err)
}
