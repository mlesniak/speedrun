package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"time"
)

// Translate all coordinates in player's X coordinate to create a fake viewport.
func drawRect(dst *ebiten.Image, x, y, w, height float64, clr color.Color) {
	translatedX := x - float64(player.Body.X) + getXTranslation()
	ebitenutil.DrawRect(dst, translatedX, y, w, height, clr)
}

func getXTranslation() float64 {
	if player.Body.X < width/2 {
		return float64(player.Body.X)
	}
	if player.Body.X >= width*widthFactor-width/2 {
		return width/2 + (float64(player.Body.X) - (width*widthFactor - width/2))
	}
	return width / 2
}

func drawHUD(screen *ebiten.Image) {
	step := int64(750)
	duration := int64(step * 4)
	passedTime := duration - time.Now().Sub(startTime).Milliseconds()
	if passedTime < step {
		hud = false
		startTime = time.Now()
		return
	}
	passedTime = passedTime / step
	secs := fmt.Sprintf("%d", int(passedTime))
	text.Draw(screen, secs, Font(160), width/2-len(secs)*50/2, height/2, color.Gray{Y: 200})
	text.Draw(screen, randomSeed.Code, Font(20), (width-len(randomSeed.Code)*10)/2, height/2+50, color.Gray{Y: 180})
}

func drawState(screen *ebiten.Image) {
	// TODO Add scenese.
	if hud {
		drawHUD(screen)
	} else {
		goal.Draw(screen)
		player.Draw(screen)
		obstacles.Draw(screen)
		levelcode.Draw(screen)
		timer.Draw(screen)
	}

	drawDebugInfo(screen)
}
