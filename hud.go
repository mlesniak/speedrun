package main

// TODO I don't like that in this file Hud is an Object and a scene.

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"time"
)

// Define HudScene
// Should the scene have an arbitrary state object, too?
var hudScene = &Scene{
	Init: func() {
		hud = NewHud()
	},
	Reset:  hud.Reset,
	Update: hud.Update,
	Draw:   hud.Draw,
}

type Hud struct {
	startTime time.Time
}

var hud = NewHud()

func NewHud() *Hud {
	return &Hud{
		startTime: time.Now(),
	}
}

func (*Hud) Update() {
	CheckExitKey()
}

func (h *Hud) Draw(screen *ebiten.Image) {
	background.Draw(screen)

	step := int64(750)
	duration := int64(step * 4)
	passedTime := duration - time.Now().Sub(h.startTime).Milliseconds()
	if passedTime < step {
		//showHud = false
		h.startTime = time.Now() // TODO remove this, too
		SetScene("game")
		return
	}
	passedTime = passedTime / step
	secs := fmt.Sprintf("%d", int(passedTime))
	text.Draw(screen, secs, Font(160), width/2-len(secs)*50/2, height/2, color.Gray{Y: 200})
	text.Draw(screen, randomSeed.Code, Font(20), (width-len(randomSeed.Code)*10)/2, height/2+50, color.Gray{Y: 180})
}

func (h *Hud) Reset() {
	hud = NewHud()
}
