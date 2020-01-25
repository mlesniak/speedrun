package main

// Pack all assets into pkged.go for single file distribution.
//
//go:generate go get github.com/markbates/pkger/cmd/pkger
//go:generate pkger -include /assets

import (
	"github.com/hajimehoshi/ebiten"
	"log"
)

func main() {
	initializeNewGame()
	if err := ebiten.Run(update, width, height, 1.0, title); err != nil {
		log.Fatal("Unable to start game: ", err)
	}
}

// update is the main game loop, updating the current game updateState and (optionally) drawing it.
func update(screen *ebiten.Image) error {
	background.Update()
	// Legacy objects.
	updateState()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	background.Draw(screen)
	// Legacy objects.
	drawState(screen)
	return nil
}
