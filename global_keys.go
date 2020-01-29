package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"os"
)

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
