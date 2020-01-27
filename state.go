// TODO Remove this class after everything is in world state?
package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"os"
)

var frameCounter = 0
var paused = false // True if the game is paused: state is not updated, but still drawn.

func CheckGameKeys() {
	if inpututil.IsKeyJustReleased(ebiten.KeyR) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton7) {
		resetGame()
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyN) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton6) {
		showHud = true // TODO Change scene...
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

	frameCounter++ // Will be removed, since each component should hold its own local state.
	if !showHud {
		player.Update()
	}
}
