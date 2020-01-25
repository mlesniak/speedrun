package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type debugFunc func() string

var debugFunctions = []debugFunc{}

// If true, debug information is displayed.
var showDebug = false

func init() {
	AddDebugMessage(func() string {
		return fmt.Sprintf("Levelcode %s", randomSeed.Code)
	})
	AddDebugMessage(func() string {
		return fmt.Sprintf("X=%.2v", player.Body.X)
	})
	AddDebugMessage(func() string {
		return fmt.Sprintf("Translation=%.2f", -float64(player.Body.X)-width/2)
	})
}

func AddDebugMessage(f debugFunc) {
	debugFunctions = append(debugFunctions, f)
}

func checkDebugKey() {
	if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		showDebug = !showDebug
	}
}

func drawDebugInfo(screen *ebiten.Image) {
	if !showDebug {
		return
	}

	rowHeight := 14
	y := 50 // Draw under level code display.
	for _, function := range debugFunctions {
		msg := function()
		ebitenutil.DebugPrintAt(screen, msg, 10, y)
		y += rowHeight
	}
}
