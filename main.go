package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
	"log"
	"os"
)

const gravity = 10

type Vector2 struct {
	X, Y float64
}

type Object struct {
	gray uint8

	X, Y     float64
	W, H     float64
	Velocity Vector2
}

var player = Object{
	gray:     40,
	X:        width / 2,
	Y:        height / 2,
	W:        20,
	H:        20,
	Velocity: Vector2{0, 0.0},
}

func main() {
	addDebugMessage(func() string {
		return fmt.Sprintf("TPS %d", ebiten.MaxTPS())
	})
	addDebugMessage(func() string {
		return fmt.Sprintf("Player.Velocity.Y=%.2f", player.Velocity.Y)
	})

	if err := ebiten.Run(update, width, height, 1, title); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	checkExitKey()

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		player.Velocity.Y -= -20
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		if player.Velocity.Y < -6.0 {
			player.Velocity.Y = -0.6
		}
	}

	// Basic physics.
	delta := 1.0 / float64(ebiten.MaxTPS())
	player.X += player.Velocity.X * delta
	player.Y += player.Velocity.Y * delta
	player.Velocity.Y -= gravity * delta

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	drawBackground(screen)
	debugInfo(screen)
	draw(screen, player)
	return nil
}

func checkExitKey() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
}

func drawBackground(screen *ebiten.Image) {
	screen.Fill(color.Gray{Y: 100})
}

func draw(screen *ebiten.Image, object Object) {
	ebitenutil.DrawRect(screen, object.X, height-object.Y, object.W, object.H, color.Gray{Y: object.gray})
}
