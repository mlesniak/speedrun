package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
)

type debugFunc func() string

// Encapsulates all debugging information.
type Debug struct {
	showDebug      bool
	debugFunctions []debugFunc
}

var debug = new(Debug)

func init() {
	AddDebugMessage(func() string {
		return fmt.Sprintf("Levelcode %s", randomSeed.Code)
	})
	AddDebugMessage(func() string {
		return fmt.Sprintf("X=%.2v", gameState.player.Body.X)
	})
	AddDebugMessage(func() string {
		return fmt.Sprintf("Translation=%.2f", -float64(gameState.player.Body.X)-width/2)
	})
}

func AddDebugMessage(f debugFunc) {
	debug.debugFunctions = append(debug.debugFunctions, f)
}

func CheckDebugKey() {
	if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		debug.showDebug = !debug.showDebug
	}
}

func (*Debug) Draw(screen *ebiten.Image) {
	if !debug.showDebug {
		return
	}

	// Show text
	rowHeight := 14
	y := 50 // Draw under level code display.
	for _, function := range debug.debugFunctions {
		msg := function()
		ebitenutil.DebugPrintAt(screen, msg, 10, y)
		y += rowHeight
	}

	// Show cursor position.
	px, py := ebiten.CursorPosition()
	crossColor := color.RGBA{80, 80, 80, 255}
	ebitenutil.DrawLine(screen, float64(px), 0, float64(px), float64(height), crossColor)
	ebitenutil.DrawLine(screen, 0, float64(py), float64(width), float64(py), crossColor)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v/%v", px, py), px+5, py+10)
}
