package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
)

// Draw a single rectangle with translated X coordinate.
func drawRect(dst *ebiten.Image, x, y, w, height float64, clr color.Color) {
	translatedX := x - float64(gameState.player.Body.X) + getXTranslation()
	ebitenutil.DrawRect(dst, translatedX, y, w, height, clr)
}

// Translate all coordinates in player's X coordinate to create a fake viewport.
func getXTranslation() float64 {
	if gameState.player.Body.X < width/2 {
		return float64(gameState.player.Body.X)
	}
	if gameState.player.Body.X >= width*widthFactor-width/2 {
		return width/2 + (float64(gameState.player.Body.X) - (width*widthFactor - width/2))
	}
	return width / 2
}
