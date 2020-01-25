package main

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
)

// Vector2 is used for physical properties as well as points. Still not sure
// if a separate Point with ints would be a viable alternative here?
type Vector2 struct {
	X, Y float64
}

// The general game object, containing of its position (Body), velocity and acceleration.
// These values are modified and used by the physics system.
type Object struct {
	Body         *resolv.Rectangle
	Velocity     Vector2
	Acceleration Vector2

	// Rename this to a color, later a texture.
	gray uint8

	// TODO Move these fields to a dedicated player object.
	// Number of times jumped
	jumped int

	// Last N positions; we could remember the timestamp, too for more independence of the frameCounter counter.
	PreviousPosition []Vector2
}

type Objecter interface {
	Update()
	Draw(*ebiten.Image) // Will this work with translations for viewports?
}
