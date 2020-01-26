package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"math"
	"time"
)

var timer *Timer = new(Timer)

type Timer struct {
}

func (*Timer) Update() {
	// Empty?
}

func (*Timer) Draw(screen *ebiten.Image) {
	var passedTime float64
	if finalTime != 0.0 {
		passedTime = finalTime
	} else {
		passedTime = time.Now().Sub(startTime).Seconds()
	}
	secs := fmt.Sprintf("%.3f", passedTime)
	text.Draw(screen, secs, Font(40), width-len(secs)*30, 45, color.Gray{Y: 200})

	if bestTime != math.MaxFloat64 {
		best := fmt.Sprintf("HIGH %.3f ", bestTime)
		text.Draw(screen, best, Font(20), width-len(best)*15, 80, color.Gray{Y: 150})
	}
}
