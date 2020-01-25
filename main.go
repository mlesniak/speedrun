package main

//go:generate go get github.com/markbates/pkger/cmd/pkger
//go:generate pkger -include /assets

import (
	"fmt"
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/mlesniak/speedrun/internal/seed"
	"log"
	"math"
	"time"
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

	// Last N positions; we could remember the timestamp, too for more independence of the frameCounter counter.
	PreviousPosition []Vector2
}

var player Object
var walls *resolv.Space
var goals *resolv.Space
var blocks = []Object{}
var goal Object

var randomSeed seed.Seed

var startTime time.Time
var finalTime = 0.0
var bestTime = math.MaxFloat64
var borderWidth int32 = 5

// Add scenes instead of a single boolean variable.
var hud = true
var countDown time.Time

func main() {
	addDebugMessage(func() string {
		return fmt.Sprintf("Levelcode %s", randomSeed.Code)
	})
	addDebugMessage(func() string {
		return fmt.Sprintf("X=%.2v", player.Body.X)
	})
	addDebugMessage(func() string {
		return fmt.Sprintf("Translation=%.2f", -float64(player.Body.X)-width/2)
	})

	// Compute scaling factor for fullscreen.
	scale := 1.0
	if fullscreen {
		ebiten.SetFullscreen(true)
		ebiten.SetCursorVisible(false)
	}

	newGame()
	if err := ebiten.Run(update, width, height, scale, title); err != nil {
		log.Fatal(err)
	}
}

var frameCounter = 0

var pause = false

func update(screen *ebiten.Image) error {
	state()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	draw(screen)
	return nil
}
