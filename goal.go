package main

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"math/rand"
)

var goals *resolv.Space
var goal Goal

type Goal struct {
	Object
}

func (g *Goal) Draw(screen *ebiten.Image) {
	drawRect(screen,
		float64(goal.Body.X), float64(goal.Body.Y), float64(goal.Body.W), float64(goal.Body.H),
		color.Gray{Y: goal.Gray})
}

func initGoals() {
	goals = resolv.NewSpace()
	x := rand.Intn(width/2) + ((width-1)*widthFactor - width/2)
	y := rand.Intn(height)
	goal = Goal{
		Object: Object{
			Gray: 255,
			Body: resolv.NewRectangle(int32(x), int32(y), player.Body.W, player.Body.H),
		}}
	goals.Add(goal.Body)
}
