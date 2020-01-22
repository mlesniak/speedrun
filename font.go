package main

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"io/ioutil"
)

var arcadeFont font.Face

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
}
