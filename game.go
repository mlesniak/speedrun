package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mlesniak/speedrun/internal/seed"
	"math"
	"math/rand"
	"os"
	"time"
)

// State
var player Player
var randomSeed seed.Seed

// Variables
var startTime time.Time

// Add scenes instead of a single boolean variable.
var showHud = true

func initalizeGame() {
	if fullscreen {
		ebiten.SetFullscreen(true)
		ebiten.SetCursorVisible(false)
	}

	randomSeed = seed.New()
	if len(os.Args) > 1 {
		randomSeed = seed.NewPreset(os.Args[1])
	}

	timer = NewTimer()

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

	timer.Reset()

	startTime = time.Now()
}
