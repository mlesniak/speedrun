package main

import (
	"fmt"
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"speedrun/seed"
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

// Add scenes instead of a single boolean variable.
var hud = true
var countDown time.Time

func main() {
	addDebugMessage(func() string {
		return fmt.Sprintf("Levelcode %s", randomSeed.Code)
	})
	addDebugMessage(func() string {
		return fmt.Sprintf("FinalTime %g", finalTime)
	})
	addDebugMessage(func() string {
		return fmt.Sprintf("Player.Velocity.Y=%.2f", player.Velocity.Y)
	})
	addDebugMessage(func() string {
		return fmt.Sprintf("Player.Velocity.X=%.2f", player.Velocity.X)
	})

	newGame()
	if err := ebiten.Run(update, width, height, 1, title); err != nil {
		log.Fatal(err)
	}
}

func newGame() {
	randomSeed = seed.New()
	initGame()
}

func initGame() {
	player = Object{
		gray:         40,
		Body:         resolv.NewRectangle(0, height-20, 20, 20),
		Velocity:     Vector2{},
		Acceleration: Vector2{X: 10.0, Y: 10.0},
	}

	// Add floor and ceiling to global space.
	walls = resolv.NewSpace()
	walls.Add(resolv.NewRectangle(0, height, width, height))
	walls.Add(resolv.NewRectangle(0, 0, width, 0))
	walls.Add(resolv.NewLine(0, 0, 0, height))
	walls.Add(resolv.NewLine(width, 0, width, height)) // Temporary, until viewport scrolling is implemented.

	// Add final goal
	goals = resolv.NewSpace()
	goal = Object{
		gray: 255,
		Body: resolv.NewRectangle(width-20-1, height-20, player.Body.W, player.Body.H),
	}
	goals.Add(goal.Body)

	// Randomize...
	rand.Seed(randomSeed.Seed)

	// Add dynamic blocks.
	blocks = []Object{}
	numBlocks := 20
	for i := 0; i < numBlocks; i++ {
		blocks = append(blocks, Object{
			gray: uint8(10 + rand.Intn(50)),
			Body: resolv.NewRectangle(rand.Int31n(width)/40*40, rand.Int31n(height)/40*40, 40, 40),
		})
	}
	for _, block := range blocks {
		walls.Add(block.Body)
	}

	// Start startTime
	finalTime = 0.0
	startTime = time.Now()

	// For HUD.
	countDown = time.Now()
}

var frameCounter = 0

func update(screen *ebiten.Image) error {
	frameCounter++
	checkExitKey()
	checkDebugKey()
	checkGameControlKeys()

	if !hud {
		updateState()
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	drawBackground(screen)

	if !hud {
		drawGoal(screen, goal)
		draw(screen, player)
		drawBlocks(screen)
		drawLevelCode(screen)
		drawTimer(screen)
	} else {
		drawHUD(screen)
	}

	debugInfo(screen)
	return nil
}

func drawHUD(screen *ebiten.Image) {
	passedTime := 4 - time.Now().Sub(startTime).Seconds()
	if passedTime < 1 {
		hud = false
		startTime = time.Now()
		return
	}
	secs := fmt.Sprintf("%d ...", int(passedTime))
	text.Draw(screen, secs, arcadeFontBig, width/2-len(secs)*30/2, height/2, color.Gray{Y: 200})
}

func updateState() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		player.jumped++
		switch player.jumped {
		case 1:
			player.Velocity.Y -= 75
		case 2:
			player.Velocity.Y -= 50
		}
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
	// Will be changed if we add viewports.
	if finalTime == 0.0 && player.Body.X == width-20-1 && player.Body.Y == height-20 {
		finalTime = time.Now().Sub(startTime).Seconds()
	}
}

func checkGameControlKeys() {
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		hud = true
		initGame()
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyN) {
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
}

func drawLevelCode(screen *ebiten.Image) {
	// Currently hard-coded, although we could use the font to retrieve the actual width and align correctly.
	text.Draw(screen, randomSeed.Code, arcadeFont, 10, 30, color.Gray{Y: 150})
}

func drawGoal(screen *ebiten.Image, object Object) {
	ebitenutil.DrawRect(screen,
		float64(object.Body.X), float64(object.Body.Y), float64(object.Body.W), float64(object.Body.H),
		color.Gray{Y: object.gray})
}

func drawBlocks(screen *ebiten.Image) {
	for _, object := range blocks {
		ebitenutil.DrawRect(screen,
			float64(object.Body.X), float64(object.Body.Y), float64(object.Body.W), float64(object.Body.H),
			color.Gray{Y: object.gray})
	}
}

func checkExitKey() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
}

func drawBackground(screen *ebiten.Image) {
	screen.Fill(color.Gray{Y: 100})
}

func draw(screen *ebiten.Image, object Object) {
	if len(object.PreviousPosition) > 0 {
		for _, vec := range object.PreviousPosition {
			ebitenutil.DrawRect(screen,
				vec.X+float64(object.Body.W)/4, vec.Y+float64(object.Body.H)/4, float64(object.Body.W)/4, float64(object.Body.H)/4,
				color.Gray{Y: object.gray * 3})
		}
	}

	ebitenutil.DrawRect(screen,
		float64(object.Body.X), float64(object.Body.Y), float64(object.Body.W), float64(object.Body.H),
		color.Gray{Y: object.gray})
}
