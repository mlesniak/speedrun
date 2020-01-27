package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"math"
	"time"
)

var timer *Timer

type Timer struct {
	highscore float64
}

func NewTimer() *Timer {
	return &Timer{
		highscore: math.MaxFloat64,
	}
}

func (*Timer) Update() {
	// Empty?
}

func (t *Timer) Draw(screen *ebiten.Image) {
	var passedTime float64
	if finalTime != 0.0 {
		passedTime = finalTime
	} else {
		passedTime = time.Now().Sub(startTime).Seconds()
	}
	secs := fmt.Sprintf("%.3f", passedTime)
	text.Draw(screen, secs, Font(40), width-len(secs)*30, 45, color.Gray{Y: 200})

	if t.highscore != math.MaxFloat64 {
		best := fmt.Sprintf("HIGH %.3f ", t.highscore)
		text.Draw(screen, best, Font(20), width-len(best)*15, 80, color.Gray{Y: 150})
	}
}

func (t *Timer) UpdateHighscore(score float64) bool {
	if score < t.highscore {
		t.highscore = score
		return true
	}

	return false
}
