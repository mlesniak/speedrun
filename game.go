package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mlesniak/speedrun/internal/seed"
	"math"
	"math/rand"
	"os"
)

var gameScene = &Scene{
	Init:   initalizeGame,
	Reset:  resetGame,
	Update: updateState,
	Draw:   drawState,
}

// State TODO Later instead of global variables since each state can carry it's own state. Global variables for shared objects? Or copying?
type GameState struct {
	player    *Player
	obstacles *Obstacles
	goal      *Goal
	timer     *Timer
}

// TODO Add functions to Game instead of global ones. So we can store state. In particular, Update and Draw do not need to have parmaeters with passed game state

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

	// TODO How to handle random seeds? Bug: hud shows other level code than later game since it's initialized two times
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
