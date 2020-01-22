package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"log"
	"os"
)

type Object struct {
	x, y float64
	gray uint8
}

var player = Object{
	x:    0,
	y:    240,
	gray: 40,
}

var obstacle = Object{
	x:    600 - 50,
	y:    240,
	gray: 80,
}

func main() {
	if err := ebiten.Run(update, 600, 480, 1, "Speedrun"); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	dx := 10.0
	player.x += dx

	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	drawBackground(screen)
	draw(screen, player)
	draw(screen, obstacle)
	return nil
}

func drawBackground(screen *ebiten.Image) {
	screen.Fill(color.Gray{Y: 100})
}

func draw(screen *ebiten.Image, object Object) {
	ebitenutil.DrawRect(screen, object.x, object.y, 50, 50, color.Gray{Y: object.gray})
}
