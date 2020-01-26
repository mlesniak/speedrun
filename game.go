package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mlesniak/speedrun/internal/seed"
	"math"
	"math/rand"
	"os"
	"time"
)

var player Player

var randomSeed seed.Seed

var startTime time.Time
var finalTime = 0.0
var bestTime = math.MaxFloat64
var borderWidth int32 = 5

// Add scenes instead of a single boolean variable.
var showHud = true
var countDown time.Time

func startGame() {
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
	resetGame()

	// For local development.
	showHud = false
}

func resetGame() {
	rand.Seed(randomSeed.Seed)

	initPlayer()
	initObstacles()
	initGoals()

	// Start startTime
	finalTime = 0.0
	startTime = time.Now()

	// For HUD.
	countDown = time.Now()
}
