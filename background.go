package main

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

var background = Background{
	color: color.Gray{Y: 100},
}

type Background struct {
	color color.Color
}

func (b *Background) Draw(screen *ebiten.Image) {
	_ = screen.Fill(b.color)
}
