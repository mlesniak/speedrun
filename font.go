package main

import (
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
)

var arcadeFont font.Face

func init() {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		panic(err)
	}

	const dpi = 72
	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
