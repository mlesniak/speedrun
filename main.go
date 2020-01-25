package main

//go:generate go get github.com/markbates/pkger/cmd/pkger
//go:generate pkger -include /assets

import (
	"fmt"
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/mlesniak/speedrun/internal/seed"
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
