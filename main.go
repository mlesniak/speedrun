package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"log"
	"os"
)

type Player struct {
	x, y float64
}

var player = Player{
	x: 0,
	y: 240,
}

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
	drawPlayer(screen)
	return nil
}

func drawPlayer(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, player.x, player.y, 50, 50, color.Gray{20})
}
