package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
)

type Levelcode struct {
	// TODO Should we add the randomSeed here?
}

func (*Levelcode) Update() {
	// Empty?
}

var levelcode = new(Levelcode)

func (*Levelcode) Draw(screen *ebiten.Image) {
	// Currently hard-coded, although we could use the font to retrieve the actual width and align correctly.
	text.Draw(screen, randomSeed.Code, Font(20), 10, 30, color.Gray{Y: 150})
}
