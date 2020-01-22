package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type debugFunc func() string

var debugFunctions = []debugFunc{}

// If true, debug information is displayed.
var showDebug = false

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
	y := 0
	for _, function := range debugFunctions {
		msg := function()
		ebitenutil.DebugPrintAt(screen, msg, 0, y)
		y += rowHeight
	}
}
