package main

// TODO How to handle random seeds? Bug: hud shows other level code than later game since it's initialized two times

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/mlesniak/speedrun/internal/seed"
	"math"
	"math/rand"
	"os"
)

var gameScene = &Scene{
	Init:   initalizeGame,
	Reset:  gameState.resetGame,
	Update: updateState,
	Draw:   drawState,
}

type GameState struct {
	player    *Player
	obstacles *Obstacles
	goal      *Goal
	timer     *Timer
}

var gameState *GameState

func initalizeGame() {
	gameState = new(GameState)

	if fullscreen {
		ebiten.SetFullscreen(true)
		ebiten.SetCursorVisible(false)
	}

	randomSeed = seed.New()
	if len(os.Args) > 1 {
		randomSeed = seed.NewPreset(os.Args[1])
	}
	rand.Seed(randomSeed.Seed)

	PlayAudioTimes("background", math.MaxInt32)

	gameState.timer = NewTimer()
	gameState.player = NewPlayer()
	gameState.obstacles = NewObstacles()
	gameState.goal = initGoals()
}

func (g *GameState) resetGame() {
	gameState.timer.Reset()
	gameState.player = NewPlayer()
}

var frameCounter = 0
var paused = false // True if the game is paused: state is not updated, but still drawn.

func CheckGameKeys() {
	if inpututil.IsKeyJustReleased(ebiten.KeyR) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton7) {
		gameState.resetGame()
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyN) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton6) {
		initalizeGame()
	}
}

func CheckExitKey() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
}

func CheckFullscreenKey() {
	if inpututil.IsKeyJustReleased(ebiten.KeyF) {
		fs := ebiten.IsFullscreen()
		ebiten.SetFullscreen(!fs)
		ebiten.SetCursorVisible(!fs)
	}
}

func CheckPauseKey() {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		paused = !paused
	}
}

func updateState() {
	CheckExitKey()
	CheckDebugKey()
	CheckFullscreenKey()
	CheckGameKeys()

	CheckPauseKey()
	if paused {
		return
	}

	frameCounter++
	gameState.player.Update()
}
