package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
)

// Translate all coordinates in player's X coordinate to create a fake viewport.
func drawRect(dst *ebiten.Image, x, y, w, height float64, clr color.Color) {
	translatedX := x - float64(gameState.player.Body.X) + getXTranslation()
	ebitenutil.DrawRect(dst, translatedX, y, w, height, clr)
}

func getXTranslation() float64 {
	if gameState.player.Body.X < width/2 {
		return float64(gameState.player.Body.X)
	}
	if gameState.player.Body.X >= width*widthFactor-width/2 {
		return width/2 + (float64(gameState.player.Body.X) - (width*widthFactor - width/2))
	}
	return width / 2
}

func drawState(screen *ebiten.Image) {
	background.Draw(screen)
	goal.Draw(screen)
	gameState.player.Draw(screen)
	gameState.obstacles.Draw(screen)
	levelcode.Draw(screen)
	gameState.timer.Draw(screen)
	debug.Draw(screen)
}
