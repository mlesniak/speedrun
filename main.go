package main

//go:generate go get github.com/markbates/pkger/cmd/pkger
//go:generate pkger -include /assets

import (
	"fmt"
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/mlesniak/speedrun/internal/seed"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

const gravity = 100

type Vector2 struct {
	X, Y float64
}

type Object struct {
	gray uint8

	Body         *resolv.Rectangle
	Velocity     Vector2
	Acceleration Vector2

	// Number of times jumped
	jumped int

	// Last N positions; we could remember the timestamp, too for more independence of the frameCounter counter.
	PreviousPosition []Vector2
}

var player Object
var walls *resolv.Space
var goals *resolv.Space
var blocks = []Object{}
var goal Object

var randomSeed seed.Seed

var startTime time.Time
var finalTime = 0.0
var bestTime = math.MaxFloat64
var borderWidth int32 = 5

// Add scenes instead of a single boolean variable.
var hud = true
var countDown time.Time

func main() {
	addDebugMessage(func() string {
		return fmt.Sprintf("Levelcode %s", randomSeed.Code)
	})
	addDebugMessage(func() string {
		return fmt.Sprintf("X=%.2v", player.Body.X)
	})
	addDebugMessage(func() string {
		return fmt.Sprintf("Translation=%.2f", -float64(player.Body.X)-width/2)
	})

	// Compute scaling factor for fullscreen.
	scale := 1.0
	if fullscreen {
		ebiten.SetFullscreen(true)
		ebiten.SetCursorVisible(false)
	}

	newGame()
	if err := ebiten.Run(update, width, height, scale, title); err != nil {
		log.Fatal(err)
	}
}

func newGame() {
	randomSeed = seed.New()
	if len(os.Args) > 1 {
		randomSeed = seed.NewPreset(os.Args[1])
	}

	bestTime = math.MaxFloat64
	playBackgroundTimes("background", math.MaxInt32)
	initGame()

	// For local development.
	hud = false
}

func initGame() {
	player = Object{
		gray:         40,
		Body:         resolv.NewRectangle(0, height-20-borderWidth, 20, 20),
		Velocity:     Vector2{},
		Acceleration: Vector2{X: 10.0, Y: 10.0},
	}

	// Add floor and ceiling to global space.
	walls = resolv.NewSpace()
	walls.Add(resolv.NewRectangle(0, height-borderWidth, (width * widthFactor), height-borderWidth))
	walls.Add(resolv.NewRectangle(0, 0+borderWidth, (width * widthFactor), 0+borderWidth))
	walls.Add(resolv.NewLine(0, 0, 0, height))
	walls.Add(resolv.NewLine((width * widthFactor), 0, (width * widthFactor), height)) // Temporary, until viewport scrolling is implemented.

	// Randomize...
	rand.Seed(randomSeed.Seed)

	// Add dynamic blocks.
	blocks = []Object{}
	for i := 0; i < numBlocks; i++ {
		blocks = append(blocks, Object{
			gray: uint8(10 + rand.Intn(50)),
			Body: resolv.NewRectangle(rand.Int31n(width*widthFactor)/40*40, rand.Int31n(height)/40*40, 40, 40),
			// TODO Check that body does not collide with player.
		})
	}
	for _, block := range blocks {
		walls.Add(block.Body)
	}

	// Add final goal in the last page of the view.
	goals = resolv.NewSpace()
	x := rand.Intn(width/2) + ((width-1)*widthFactor - width/2)
	y := rand.Intn(height)
	goal = Object{
		gray: 255,
		Body: resolv.NewRectangle(int32(x), int32(y), player.Body.W, player.Body.H),
	}
	goals.Add(goal.Body)

	// Start startTime
	finalTime = 0.0
	startTime = time.Now()

	// For HUD.
	countDown = time.Now()
}

var frameCounter = 0

var pause = false

func update(screen *ebiten.Image) error {
	checkExitKey()
	if showDebug && inpututil.IsKeyJustPressed(ebiten.KeyP) {
		pause = !pause
	}
	if !pause {
		frameCounter++
		checkExitKey()
		checkDebugKey()
		checkFullscreenKey()
		checkGameControlKeys()

		if !hud {
			updateState()
		}
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	drawBackground(screen)

	if !hud {
		drawGoal(screen, goal)
		drawPlayer(screen, player)
		drawBlocks(screen)
		drawBorders(screen)
		drawLevelCode(screen)
		drawTimer(screen)
	} else {
		drawHUD(screen)
	}

	debugInfo(screen)
	debugDrawPosition(screen)
	return nil
}

func drawBorders(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, width*widthFactor, float64(borderWidth), color.Gray{Y: 40})
	ebitenutil.DrawRect(screen, 0, float64(height-borderWidth), width*widthFactor, float64(height-borderWidth), color.Gray{Y: 40})
}

// TODO move to audio.go
func drawHUD(screen *ebiten.Image) {
	step := int64(750)
	duration := int64(step * 4)
	passedTime := duration - time.Now().Sub(startTime).Milliseconds()
	if passedTime < step {
		hud = false
		startTime = time.Now()
		return
	}
	passedTime = passedTime / step
	secs := fmt.Sprintf("%d", int(passedTime))
	text.Draw(screen, secs, arcadeFontLarge, width/2-len(secs)*50/2, height/2, color.Gray{Y: 200})
	text.Draw(screen, randomSeed.Code, arcadeFont, (width-len(randomSeed.Code)*10)/2, height/2+50, color.Gray{Y: 180})
}

func updateState() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton0) {
		player.jumped++
		switch player.jumped {
		case 1:
			player.Velocity.Y -= float64(gravity) * 0.75
		case 2:
			player.Velocity.Y -= float64(gravity) * 0.50
		}
	}

	// Check axes of first gamepad.
	axes := []float64{
		ebiten.GamepadAxis(0, 0),
		ebiten.GamepadAxis(0, 1)}
	gamepadAcceleration := 40.0
	if math.Abs(axes[0]) > 0.15 {
		player.Velocity.X += axes[0] * gamepadAcceleration
	}
	if math.Abs(axes[1]) > 0.15 {
		player.Velocity.X -= axes[1] * gamepadAcceleration
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		player.Velocity.X -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		player.Velocity.X += 5
	}
	for math.Abs(player.Velocity.X) > 40 {
		d := -0.01
		if math.Signbit(player.Velocity.X) {
			d *= -1
		}
		player.Velocity.X += d
	}

	// Basic physics.
	delta := 1.0 / float64(ebiten.MaxTPS())
	dx := int32(player.Velocity.X * delta * player.Acceleration.X)
	dy := int32(player.Velocity.Y * delta * player.Acceleration.Y)

	collision := walls.Resolve(player.Body, dx, 0)
	if collision.Colliding() {
		dx = collision.ResolveX
	}
	player.Body.X += dx
	if collision.Colliding() {
		player.Velocity.X = 0
	} else {
		if player.jumped > 0 {
			player.Velocity.X *= 0.99
		} else {
			player.Velocity.X *= 0.9
		}
		if math.Abs(player.Velocity.X) < 0.01 && player.jumped == 0 {
			player.Velocity.X = 0
		}
	}

	collision = walls.Resolve(player.Body, 0, dy)
	if collision.Colliding() {
		dy = collision.ResolveY
	}
	player.Body.Y += dy
	if collision.Colliding() {
		player.Velocity.Y = 0
		player.jumped = 0
	} else {
		player.Velocity.Y += gravity * delta
	}

	// Update historic positions.
	if frameCounter%numHistoricFramesUpdate == 0 {
		if len(player.PreviousPosition) < numHistoricFrames {
			player.PreviousPosition = append(player.PreviousPosition, Vector2{float64(player.Body.X), float64(player.Body.Y)})
		} else {
			player.PreviousPosition = player.PreviousPosition[1:]
			player.PreviousPosition = append(player.PreviousPosition, Vector2{float64(player.Body.X), float64(player.Body.Y)})
		}
	}

	// Check if goal reached.
	if finalTime == 0.0 && goals.IsColliding(player.Body) {
		finalTime = time.Now().Sub(startTime).Seconds()
		if finalTime < bestTime {
			bestTime = finalTime
			playBackgroundTimes("goal", 2)
		} else {
			playBackground("goal")
		}
	}
}

func checkGameControlKeys() {
	if inpututil.IsKeyJustReleased(ebiten.KeyR) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton7) {
		initGame()
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyN) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton6) {
		hud = true
		newGame()
	}
}

func drawTimer(screen *ebiten.Image) {
	var passedTime float64
	if finalTime != 0.0 {
		passedTime = finalTime
	} else {
		passedTime = time.Now().Sub(startTime).Seconds()
	}
	secs := fmt.Sprintf("%.3f", passedTime)
	text.Draw(screen, secs, arcadeFontBig, width-len(secs)*30, 45, color.Gray{Y: 200})

	if bestTime != math.MaxFloat64 {
		best := fmt.Sprintf("HIGH %.3f ", bestTime)
		text.Draw(screen, best, arcadeFont, width-len(best)*15, 80, color.Gray{Y: 150})
	}
}

func drawLevelCode(screen *ebiten.Image) {
	// Currently hard-coded, although we could use the font to retrieve the actual width and align correctly.
	text.Draw(screen, randomSeed.Code, arcadeFont, 10, 30, color.Gray{Y: 150})
}

func drawGoal(screen *ebiten.Image, object Object) {
	drawRect(screen,
		float64(object.Body.X), float64(object.Body.Y), float64(object.Body.W), float64(object.Body.H),
		color.Gray{Y: object.gray})
}

func drawBlocks(screen *ebiten.Image) {
	for _, object := range blocks {
		drawRect(screen,
			float64(object.Body.X), float64(object.Body.Y), float64(object.Body.W), float64(object.Body.H),
			color.Gray{Y: object.gray})
	}
}

func checkExitKey() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
}

func checkFullscreenKey() {
	if inpututil.IsKeyJustReleased(ebiten.KeyF) {
		fs := ebiten.IsFullscreen()
		ebiten.SetFullscreen(!fs)
		ebiten.SetCursorVisible(!fs)
	}
}

var backgroundColor = color.Gray{Y: 100}

func drawBackground(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
}

func drawPlayer(screen *ebiten.Image, object Object) {
	x := getXTranslation()

	// Trail
	if len(object.PreviousPosition) > 0 {
		colorQuotient := float64(backgroundColor.Y) / float64(len(object.PreviousPosition))
		for i, vec := range object.PreviousPosition {
			boxColor := uint8(float64(backgroundColor.Y) + colorQuotient*float64(i))
			drawRect(screen,
				vec.X+float64(object.Body.W)/2-float64(object.Body.W)/8,
				vec.Y+float64(object.Body.H)/2-float64(object.Body.H)/8,
				float64(object.Body.W)/4,
				float64(object.Body.H)/4,
				color.Gray{Y: boxColor})
		}
	}

	col := color.RGBA{255, 255, 0, 255}
	ebitenutil.DrawRect(screen, x, float64(object.Body.Y), float64(object.Body.W), float64(object.Body.H), col)
}

func getXTranslation() float64 {
	if player.Body.X < width/2 {
		return float64(player.Body.X)
	}
	if player.Body.X >= width*widthFactor-width/2 {
		return width/2 + (float64(player.Body.X) - (width*widthFactor - width/2))
	}
	return width / 2
}

// Translate all coordinates in player's X coordinate to create a fake viewport.
func drawRect(dst *ebiten.Image, x, y, w, height float64, clr color.Color) {
	translatedX := x - float64(player.Body.X) + getXTranslation()
	ebitenutil.DrawRect(dst, translatedX, y, w, height, clr)
}

func debugDrawPosition(screen *ebiten.Image) {
	if !showDebug {
		return
	}

	px, py := ebiten.CursorPosition()
	crossColor := color.RGBA{80, 80, 80, 255}
	ebitenutil.DrawLine(screen, float64(px), 0, float64(px), float64(height), crossColor)
	ebitenutil.DrawLine(screen, 0, float64(py), float64(width), float64(py), crossColor)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v/%v", px, py), px+5, py+10)
}
