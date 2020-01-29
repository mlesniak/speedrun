package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"math"
	"time"
)

type Timer struct {
	highscore float64
	goalTime  float64
	startTime time.Time
}

func NewTimer() *Timer {
	return &Timer{
		highscore: math.MaxFloat64,
		goalTime:  math.MaxFloat64,
		startTime: time.Now(),
	}
}

func (*Timer) Update() {
	// Empty?
}

func (t *Timer) Draw(screen *ebiten.Image) {
	// Draw passed time or reached time, if available.
	var displayTime float64
	if t.goalTime != math.MaxFloat64 {
		displayTime = t.goalTime
	} else {
		displayTime = time.Now().Sub(t.startTime).Seconds()
	}
	fmtTime := fmt.Sprintf("%.3f", displayTime)
	text.Draw(screen, fmtTime, Font(40), width-len(fmtTime)*30, 45, color.Gray{Y: 200})

	// Draw highscore if available.
	if t.highscore != math.MaxFloat64 {
		fmtHighscore := fmt.Sprintf("HIGH %.3f ", t.highscore)
		text.Draw(screen, fmtHighscore, Font(20), width-len(fmtHighscore)*15, 80, color.Gray{Y: 150})
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

	now := time.Now().Sub(t.startTime).Seconds()
	t.goalTime = math.Min(now, t.goalTime)
	return t.updateHighscore(t.goalTime)
}

// Reset resets the timer but keeps the highscore.
func (t *Timer) Reset() {
	t.goalTime = math.MaxFloat64
	t.startTime = time.Now()
}
