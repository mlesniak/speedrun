package main

import (
	"fmt"
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
	"log"
	"math"
	"os"
)

const gravity = 100

type Vector2 struct {
	X, Y float64
}

type Object struct {
	gray uint8

	Body         *resolv.Rectangle
	Velocity     Vector2
	Acceleration Vector2

	// Number of times jumped
	jumped int
}

var player = Object{
	gray:         40,
	Body:         resolv.NewRectangle(width/2, height*0.8, 20, 20),
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
	addDebugMessage(func() string {
		return fmt.Sprintf("Player.Velocity.X=%.2f", player.Velocity.X)
	})

	// Add floor and ceiling to global space.
	walls = resolv.NewSpace()
	walls.Add(resolv.NewRectangle(0, height, width, height))
	walls.Add(resolv.NewRectangle(0, 0, width, 0))

	if err := ebiten.Run(update, width, height, 1, title); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	checkExitKey()

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		player.jumped++
		switch player.jumped {
		case 1:
			player.Velocity.Y -= 75
		case 2:
			player.Velocity.Y -= 50
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		player.Velocity.X -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		player.Velocity.X += 5
	}
	for math.Abs(player.Velocity.X) > 40 {
		d := -0.01
		if math.Signbit(player.Velocity.X) {
			d *= -1
		}
		player.Velocity.X += d
	}

	// Basic physics.
	delta := 1.0 / float64(ebiten.MaxTPS())
	dx := int32(player.Velocity.X * delta * player.Acceleration.X)
	dy := int32(player.Velocity.Y * delta * player.Acceleration.Y)

	collision := walls.Resolve(player.Body, dx, 0)
	if collision.Colliding() {
		dx = collision.ResolveX
	}
	player.Body.X += dx
	if collision.Colliding() {
		player.Velocity.X = 0
	} else {
		if player.jumped > 0 {
			player.Velocity.X *= 0.99
		} else {
			player.Velocity.X *= 0.9
		}
		if math.Abs(player.Velocity.X) < 0.01 && player.jumped == 0 {
			player.Velocity.X = 0
		}
	}

	collision = walls.Resolve(player.Body, 0, dy)
	if collision.Colliding() {
		dy = collision.ResolveY
	}
	player.Body.Y += dy
	if collision.Colliding() {
		player.Velocity.Y = 0
		player.jumped = 0
	} else {
		player.Velocity.Y += gravity * delta
	}

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
