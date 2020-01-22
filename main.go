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
	gray uint8
	resolv.Rectangle
}

var player = Object{
	gray:      40,
	Rectangle: *resolv.NewRectangle(0, 240, 50, 50),
}

var obstacle = Object{
	gray:      80,
	Rectangle: *resolv.NewRectangle(600-50, 240, 50, 50),
}

func main() {
	if err := ebiten.Run(update, 600, 480, 1, "Speedrun"); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	// Simulate input.
	// TODO What happens with gamepad input < 1.0?
	dx := int32(10)

	// Check for collision and update accordingly.
	collision := resolv.Resolve(&player.Rectangle, &obstacle.Rectangle, dx, 0)
	if collision.Colliding() {
		dx = collision.ResolveX
	}
	player.X += dx

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
	ebitenutil.DrawRect(screen, float64(object.X), float64(object.Y), float64(object.W), float64(object.H), color.Gray{Y: object.gray})
}
