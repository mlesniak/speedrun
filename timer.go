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
	goalTime  float64
}

func NewTimer() *Timer {
	return &Timer{
		highscore: math.MaxFloat64,
		goalTime:  math.MaxFloat64,
	}
}

func (*Timer) Update() {
	// Empty?
}

func (t *Timer) Draw(screen *ebiten.Image) {
	var passedTime float64
	if t.goalTime != math.MaxFloat64 {
		passedTime = t.goalTime
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

func (t *Timer) updateHighscore(score float64) bool {
	if score < t.highscore {
		t.highscore = score
		return true
	}

	return false
}

// UpdateTime updates the final time in the timer, if not done yet.
func (t *Timer) UpdateTime() bool {
	// Already set? Ignore update.
	if t.goalTime != math.MaxFloat64 {
		return false
	}

	now := time.Now().Sub(startTime).Seconds()
	t.goalTime = math.Min(now, t.goalTime)
	return t.updateHighscore(t.goalTime)
}

// Reset resets the timer but keeps the highscore.
func (t *Timer) Reset() {
	t.goalTime = math.MaxFloat64
}
