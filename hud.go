package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"time"
)

// Empty, but not yet fully functional. This approach is necessary, since method values bind their receiver, i.e. see my question
// on StackOverflow at https://stackoverflow.com/questions/59971005/why-is-this-behaviour-regarding-method-pointers-and-global-variable-initializati
var hudState = new(HudState)

// Define HudScene
var hudScene = &Scene{
	Init:   hudState.initHud,
	Reset:  hudState.Reset,
	Update: hudState.Update,
	Draw:   hudState.Draw,
}

type HudState struct {
	startTime time.Time
}

func (h *HudState) initHud() {
	InitializeRandomSeed()
	h.startTime = time.Now()
}

func (*HudState) Update() {
	CheckExitKey()
}

func (h *HudState) Draw(screen *ebiten.Image) {
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

func (h *HudState) Reset() {
	h.initHud()
}
