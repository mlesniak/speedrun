package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/mlesniak/speedrun/internal/seed"
	"math"
	"math/rand"
)

var gameScene = &Scene{
	Init: func() {
		initalizeGame()
	},
	Reset:  gameState.resetGame, // BUG is nil of course...
	Update: gameState.updateState,
	Draw:   drawState,
}

type GameState struct {
	player    *Player
	obstacles *Obstacles
	goal      *Goal
	timer     *Timer

	frameCounter int  // For drawing the tail.
	paused       bool // True if we should only check for keys but not change player position.
}

var gameState *GameState

func initalizeGame() {
	gameState = new(GameState)

	if fullscreen {
		ebiten.SetFullscreen(true)
		ebiten.SetCursorVisible(false)
	}

	// For debugging. Usually done by the hud.
	randomSeed = seed.New()
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

func CheckGameKeys() {
	if inpututil.IsKeyJustReleased(ebiten.KeyR) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton7) {
		gameState.resetGame()
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyN) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton6) {
		initalizeGame()
	}
}

func CheckPauseKey(paused *bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		*paused = !(*paused)
	}
}

func (g *GameState) updateState() {
	CheckDebugKey()
	CheckGameKeys()

	CheckPauseKey(&g.paused)
	if g.paused {
		return
	}

	g.frameCounter++
	g.player.Update(g)
}
