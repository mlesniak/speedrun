package main

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"io/ioutil"
)

var arcadeFont font.Face
var arcadeFontBig font.Face
var arcadeFontLarge font.Face

func init() {
	pix, err := ioutil.ReadFile("assets/arcadepix.ttf")
	if err != nil {
		panic(err)
	}

	tt, err := truetype.Parse(pix)
	if err != nil {
		panic(err)
	}

	const dpi = 72
	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    20,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	arcadeFontBig = truetype.NewFace(tt, &truetype.Options{
		Size:    40,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	arcadeFontLarge = truetype.NewFace(tt, &truetype.Options{
		Size:    80,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
