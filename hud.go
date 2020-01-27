package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"time"
)

type Hud struct {
	// Empty
}

var hud *Hud = new(Hud)

func (*Hud) Update() {
	// Empty
}

func (*Hud) Draw(screen *ebiten.Image) {
	// TODO Countdown should not use start time of the timer but should have a timer on its own.

	step := int64(750)
	duration := int64(step * 4)
	passedTime := duration - time.Now().Sub(startTime).Milliseconds()
	if passedTime < step {
		showHud = false
		startTime = time.Now()
		return
	}
	passedTime = passedTime / step
	secs := fmt.Sprintf("%d", int(passedTime))
	text.Draw(screen, secs, Font(160), width/2-len(secs)*50/2, height/2, color.Gray{Y: 200})
	text.Draw(screen, randomSeed.Code, Font(20), (width-len(randomSeed.Code)*10)/2, height/2+50, color.Gray{Y: 180})
}
