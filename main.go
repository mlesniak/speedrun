package main

//go:generate go get github.com/markbates/pkger/cmd/pkger
//go:generate pkger -include /assets

import (
	"github.com/hajimehoshi/ebiten"
	"log"
)

func main() {
	if fullscreen {
		ebiten.SetFullscreen(true)
		ebiten.SetCursorVisible(false)
	}

	newGame()
	if err := ebiten.Run(update, width, height, 1.0, title); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	state()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	draw(screen)
	return nil
}
