package main

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
	"math"
)

type Player struct {
	Object

	// Number of times jumped
	jumped int

	// Last N positions; we could remember the timestamp, too for more independence of the frameCounter counter.
	PreviousPosition []Vector2
}

const gravity = 100

func initPlayer() {
	player = Player{
		Object: Object{
			Body:         resolv.NewRectangle(0, height-20-borderWidth, 20, 20),
			Velocity:     Vector2{},
			Acceleration: Vector2{X: 10.0, Y: 10.0},
		},
		jumped:           0,
		PreviousPosition: nil,
	}
}

func (player *Player) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton0) {
		player.jumped++
		switch player.jumped {
		case 1:
			player.Velocity.Y -= float64(gravity) * 0.75
		case 2:
			player.Velocity.Y -= float64(gravity) * 0.50
		}
	}

	// Check axes of first gamepad.
	axes := []float64{
		ebiten.GamepadAxis(0, 0),
		ebiten.GamepadAxis(0, 1)}
	gamepadAcceleration := 40.0
	if math.Abs(axes[0]) > 0.15 {
		player.Velocity.X += axes[0] * gamepadAcceleration
	}
	if math.Abs(axes[1]) > 0.15 {
		player.Velocity.X -= axes[1] * gamepadAcceleration
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

	// Update historic positions.
	if frameCounter%numHistoricFramesUpdate == 0 {
		if len(player.PreviousPosition) < numHistoricFrames {
			player.PreviousPosition = append(player.PreviousPosition, Vector2{float64(player.Body.X), float64(player.Body.Y)})
		} else {
			player.PreviousPosition = player.PreviousPosition[1:]
			player.PreviousPosition = append(player.PreviousPosition, Vector2{float64(player.Body.X), float64(player.Body.Y)})
		}
	}

	// Check if goal reached.
	// TODO Check that we did not reach the goal already (as game state, i.e. later) to prevent sound on each collision.
	if goals.IsColliding(player.Body) {
		if timer.UpdateTime() {
			PlayAudioTimes("goal", 2)
		} else {
			PlayAudio("goal")
		}
	}
}

func (player *Player) Draw(screen *ebiten.Image) {
	x := getXTranslation()

	// Trail
	if len(player.PreviousPosition) > 0 {
		r, g, b, _ := color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		}.RGBA()
		colorQuotient := 255.0 / float64(len(player.PreviousPosition))
		for i, vec := range player.PreviousPosition {
			boxColor := uint8(colorQuotient * float64(i))
			drawRect(screen,
				vec.X+float64(player.Body.W)/2-float64(player.Body.W)/8,
				vec.Y+float64(player.Body.H)/2-float64(player.Body.H)/8,
				float64(player.Body.W)/4,
				float64(player.Body.H)/4,
				color.RGBA{
					R: uint8(r),
					G: uint8(g),
					B: uint8(b),
					A: boxColor,
				})
		}
	}

	col := color.RGBA{255, 255, 0, 255}
	ebitenutil.DrawRect(screen, x, float64(player.Body.Y), float64(player.Body.W), float64(player.Body.H), col)
}
