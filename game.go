package main

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/mlesniak/speedrun/internal/seed"
	"math"
	"math/rand"
	"os"
	"time"
)

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

func initializeNewGame() {
	if fullscreen {
		ebiten.SetFullscreen(true)
		ebiten.SetCursorVisible(false)
	}

	randomSeed = seed.New()
	if len(os.Args) > 1 {
		randomSeed = seed.NewPreset(os.Args[1])
	}

	bestTime = math.MaxFloat64
	PlayAudioTimes("background", math.MaxInt32)
	resetCurrentGame()

	// For local development.
	hud = false
}

func resetCurrentGame() {
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
