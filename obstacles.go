package main

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"math/rand"
)

//var obstacles *Obstacles

// Internal
var walls *resolv.Space
var blocks = []Object{}

type Obstacles struct {
	// Shold we have an object for each block / wall or should we look at all walls as one (for now)?
	//
	// Empty for now...
}

func (*Obstacles) Update() {
	// Empty.
}

func (*Obstacles) Draw(screen *ebiten.Image) {
	drawBlocks(screen)
	drawBorders(screen)
}

func NewObstacles() *Obstacles {
	// Add floor and ceiling to global space.
	walls = resolv.NewSpace()
	walls.Add(resolv.NewRectangle(0, height-borderWidth, (width * widthFactor), height-borderWidth))
	walls.Add(resolv.NewRectangle(0, 0+borderWidth, (width * widthFactor), 0+borderWidth))
	walls.Add(resolv.NewLine(0, 0, 0, height))
	walls.Add(resolv.NewLine((width * widthFactor), 0, (width * widthFactor), height)) // Temporary, until viewport scrolling is implemented.

	// Add dynamic blocks.
	blocks = []Object{}
	for i := 0; i < numBlocks; i++ {
		blocks = append(blocks, Object{
			Gray: uint8(10 + rand.Intn(50)),
			Body: resolv.NewRectangle(rand.Int31n(width*widthFactor)/40*40, rand.Int31n(height)/40*40, 40, 40),
			// TODO Check that body does not collide with player.
		})
	}
	for _, block := range blocks {
		walls.Add(block.Body)
	}

	return new(Obstacles)
}

func drawBlocks(screen *ebiten.Image) {
	for _, object := range blocks {
		drawRect(screen,
			float64(object.Body.X), float64(object.Body.Y), float64(object.Body.W), float64(object.Body.H),
			color.Gray{Y: object.Gray})
	}
}

func drawBorders(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, width*widthFactor, float64(borderWidth), color.Gray{Y: 40})
	ebitenutil.DrawRect(screen, 0, float64(height-borderWidth), width*widthFactor, float64(height-borderWidth), color.Gray{Y: 40})
}
