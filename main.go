package main

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"log"
	"os"
)

type Object struct {
	x, y  float64
	gray  uint8
	shape resolv.Rectangle
}

var player = Object{
	x:     0,
	y:     240,
	gray:  40,
	shape: *resolv.NewRectangle(0, 240, 50, 50),
}

var obstacle = Object{
	x:     600 - 50,
	y:     240,
	gray:  80,
	shape: *resolv.NewRectangle(600-50, 240, 50, 50),
}

func main() {
	if err := ebiten.Run(update, 600, 480, 1, "Speedrun"); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	dx := 20.0
	// Check for collision.
	collision := resolv.Resolve(&player.shape, &obstacle.shape, int32(dx), 0)
	if collision.Colliding() {
		dx = float64(collision.ResolveX)
	}

	player.x += dx
	player.shape.X += int32(dx)

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
