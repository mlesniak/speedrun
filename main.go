package main

import (
	"fmt"
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
	"log"
	"os"
)

const gravity = 20

type Vector2 struct {
	X, Y float64
}

type Object struct {
	gray uint8

	Body         *resolv.Rectangle
	Velocity     Vector2
	Acceleration Vector2
}

var player = Object{
	gray:         40,
	Body:         resolv.NewRectangle(width/2, height/2, 20, 20),
	Velocity:     Vector2{},
	Acceleration: Vector2{X: 10.0, Y: 10.0},
}

var walls *resolv.Space

func main() {
	addDebugMessage(func() string {
		return fmt.Sprintf("TPS %d", ebiten.MaxTPS())
	})
	addDebugMessage(func() string {
		return fmt.Sprintf("Player.Velocity.Y=%.2f", player.Velocity.Y)
	})

	// Add floor to global space.
	walls = resolv.NewSpace()
	floor := resolv.NewRectangle(0, height, width, height)
	walls.Add(floor)

	if err := ebiten.Run(update, width, height, 1, title); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	checkExitKey()

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		player.Velocity.Y -= 20
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		if player.Velocity.Y > 6.0 {
			player.Velocity.Y = 0.6
		}
	}

	// Basic physics.
	delta := 1.0 / float64(ebiten.MaxTPS())
	dx := int32(player.Velocity.X * delta * player.Acceleration.X)
	dy := int32(player.Velocity.Y * delta * player.Acceleration.Y)

	collision := walls.Resolve(player.Body, dx, 0)
	if collision.Colliding() {
		dx = collision.ResolveX
	}
	collision = walls.Resolve(player.Body, 0, dy)
	if collision.Colliding() {
		dy = collision.ResolveY
	}

	player.Body.X += dx
	player.Body.Y += dy
	player.Velocity.Y += gravity * delta

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
	ebitenutil.DrawRect(screen,
		float64(object.Body.X), float64(object.Body.Y), float64(object.Body.W), float64(object.Body.H),
		color.Gray{Y: object.gray})
}
