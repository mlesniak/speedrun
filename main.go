package main

import (
	"fmt"
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
	if err := ebiten.Run(update, width, height, 1, title); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	// Simulate input.
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
	debugInfo(screen)
	draw(screen, player)
	draw(screen, obstacle)
	return nil
}

func debugInfo(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %d", ebiten.MaxTPS()))
}

func drawBackground(screen *ebiten.Image) {
	screen.Fill(color.Gray{Y: 100})
}

func draw(screen *ebiten.Image, object Object) {
	ebitenutil.DrawRect(screen, float64(object.X), float64(object.Y), float64(object.W), float64(object.H), color.Gray{Y: object.gray})
}
