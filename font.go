package main

import (
	"github.com/golang/freetype/truetype"
	"github.com/markbates/pkger"
	"golang.org/x/image/font"
	"io/ioutil"
	"log"
)

var fontCache = make(map[float64]font.Face)

var arcadeTruetypeFont *truetype.Font

func init() {
	// Read font.
	pix, err := pkger.Open("/assets/arcadepix.ttf")
	defer pix.Close()
	mustFont(err)
	bs, err := ioutil.ReadAll(pix)
	mustFont(err)

	// Parse font data into fontFace.
	arcadeTruetypeFont, err = truetype.Parse(bs)
	mustFont(err)
}

// Font returns the arcade font for the given size. If it is not yet generated, it will be loaded and created.
func Font(size float64) font.Face {
	if f, ok := fontCache[size]; ok {
		return f
	}

	loadedFont := createFont(size)
	fontCache[size] = loadedFont
	return loadedFont
}

// createFont creates the actual font face for the given size.
func createFont(size float64) font.Face {
	const dpi = 72
	loadedFont := truetype.NewFace(arcadeTruetypeFont, &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	return loadedFont
}

// mustAudio checks if an error occured while loading a file.
func mustFont(err error) {
	if err != nil {
		log.Fatal("Unable to load font:", err)
	}
}
