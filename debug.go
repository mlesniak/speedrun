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
	addDebugMessage(func() string {
		return fmt.Sprintf("Levelcode %s", randomSeed.Code)
	})
	addDebugMessage(func() string {
		return fmt.Sprintf("X=%.2v", player.Body.X)
	})
	addDebugMessage(func() string {
		return fmt.Sprintf("Translation=%.2f", -float64(player.Body.X)-width/2)
	})
}

func addDebugMessage(f debugFunc) {
	debugFunctions = append(debugFunctions, f)
}

func checkDebugKey() {
	if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		showDebug = !showDebug
	}
}

func debugInfo(screen *ebiten.Image) {
	if !showDebug {
		return
	}

	rowHeight := 14
	y := 50
	for _, function := range debugFunctions {
		msg := function()
		ebitenutil.DebugPrintAt(screen, msg, 10, y)
		y += rowHeight
	}
}
