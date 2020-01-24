package main

import (
	"github.com/golang/freetype/truetype"
	"github.com/markbates/pkger"
	"golang.org/x/image/font"
	"io/ioutil"
)

var arcadeFont font.Face
var arcadeFontBig font.Face
var arcadeFontLarge font.Face

func init() {
	pix, err := pkger.Open("/assets/arcadepix.ttf")
	defer pix.Close()
	if err != nil {
		panic(err)
	}
	bs, err := ioutil.ReadAll(pix)
	if err != nil {
		panic(err)
	}

	tt, err := truetype.Parse(bs)
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
		Size:    160,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
