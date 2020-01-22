package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"log"
	"os"
)

func main() {
	if err := ebiten.Run(update, 600, 480, 1, "Speedrun"); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.Fill(color.Gray{80})
	ebitenutil.DrawRect(screen, 10, 10, 100, 100, color.White)
	return nil
}
