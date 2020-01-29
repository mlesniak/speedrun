package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"time"
)

var hud *Hud

// Define HudScene
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

func NewHud() *Hud {
	InitializeRandomSeed()
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
	duration := step * 4
	passedTime := duration - time.Now().Sub(h.startTime).Milliseconds()
	if passedTime < step {
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
