package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mlesniak/speedrun/internal/seed"
	"math"
	"math/rand"
	"os"
)

// State
type GameState struct {
	player    *Player
	obstacles *Obstacles
	goal      *Goal
	timer     *Timer
}

// TODO Where to add auxiliary objects such as music? Part of a scene (init and reset, ...)?

var player Player
var randomSeed seed.Seed

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

	// TODO Currently leading to a bug where timer already starts.
	timer.Reset()
	hud.Reset()
}
