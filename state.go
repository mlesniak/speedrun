package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"math"
	"os"
	"time"
)

const gravity = 100

var frameCounter = 0

var pause = false

func updateState() {
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
	if finalTime == 0.0 && goals.IsColliding(player.Body) {
		finalTime = time.Now().Sub(startTime).Seconds()
		if finalTime < bestTime {
			bestTime = finalTime
			playBackgroundTimes("goal", 2)
		} else {
			playBackground("goal")
		}
	}
}

func checkGameControlKeys() {
	if inpututil.IsKeyJustReleased(ebiten.KeyR) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton7) {
		initGame()
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyN) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton6) {
		hud = true
		newGame()
	}
}

func checkExitKey() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
}

func checkFullscreenKey() {
	if inpututil.IsKeyJustReleased(ebiten.KeyF) {
		fs := ebiten.IsFullscreen()
		ebiten.SetFullscreen(!fs)
		ebiten.SetCursorVisible(!fs)
	}
}

func state() {
	checkExitKey()
	if showDebug && inpututil.IsKeyJustPressed(ebiten.KeyP) {
		pause = !pause
	}
	if !pause {
		frameCounter++
		checkExitKey()
		checkDebugKey()
		checkFullscreenKey()
		checkGameControlKeys()

		if !hud {
			updateState()
		}
	}
}
