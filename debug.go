package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type debugFunc func() string

var debugFunctions = []debugFunc{}

func addDebugMessage(f debugFunc) {
	debugFunctions = append(debugFunctions, f)
}

func debugInfo(screen *ebiten.Image) {
	rowHeight := 14
	y := 0
	for _, function := range debugFunctions {
		msg := function()
		ebitenutil.DebugPrintAt(screen, msg, 0, y)
		y += rowHeight
	}
}
