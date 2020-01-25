package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"math"
	"time"
)

var backgroundColor = color.Gray{Y: 100}

func drawBackground(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
}

func drawPlayer(screen *ebiten.Image, object Object) {
	x := getXTranslation()

	// Trail
	if len(object.PreviousPosition) > 0 {
		colorQuotient := float64(backgroundColor.Y) / float64(len(object.PreviousPosition))
		for i, vec := range object.PreviousPosition {
			boxColor := uint8(float64(backgroundColor.Y) + colorQuotient*float64(i))
			drawRect(screen,
				vec.X+float64(object.Body.W)/2-float64(object.Body.W)/8,
				vec.Y+float64(object.Body.H)/2-float64(object.Body.H)/8,
				float64(object.Body.W)/4,
				float64(object.Body.H)/4,
				color.Gray{Y: boxColor})
		}
	}

	col := color.RGBA{255, 255, 0, 255}
	ebitenutil.DrawRect(screen, x, float64(object.Body.Y), float64(object.Body.W), float64(object.Body.H), col)
}

// Translate all coordinates in player's X coordinate to create a fake viewport.
func drawRect(dst *ebiten.Image, x, y, w, height float64, clr color.Color) {
	translatedX := x - float64(player.Body.X) + getXTranslation()
	ebitenutil.DrawRect(dst, translatedX, y, w, height, clr)
}

func debugDrawPosition(screen *ebiten.Image) {
	if !showDebug {
		return
	}

	px, py := ebiten.CursorPosition()
	crossColor := color.RGBA{80, 80, 80, 255}
	ebitenutil.DrawLine(screen, float64(px), 0, float64(px), float64(height), crossColor)
	ebitenutil.DrawLine(screen, 0, float64(py), float64(width), float64(py), crossColor)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v/%v", px, py), px+5, py+10)
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

func drawTimer(screen *ebiten.Image) {
	var passedTime float64
	if finalTime != 0.0 {
		passedTime = finalTime
	} else {
		passedTime = time.Now().Sub(startTime).Seconds()
	}
	secs := fmt.Sprintf("%.3f", passedTime)
	text.Draw(screen, secs, arcadeFontBig, width-len(secs)*30, 45, color.Gray{Y: 200})

	if bestTime != math.MaxFloat64 {
		best := fmt.Sprintf("HIGH %.3f ", bestTime)
		text.Draw(screen, best, arcadeFont, width-len(best)*15, 80, color.Gray{Y: 150})
	}
}

func drawLevelCode(screen *ebiten.Image) {
	// Currently hard-coded, although we could use the font to retrieve the actual width and align correctly.
	text.Draw(screen, randomSeed.Code, arcadeFont, 10, 30, color.Gray{Y: 150})
}

func drawGoal(screen *ebiten.Image, object Object) {
	drawRect(screen,
		float64(object.Body.X), float64(object.Body.Y), float64(object.Body.W), float64(object.Body.H),
		color.Gray{Y: object.gray})
}

func drawBlocks(screen *ebiten.Image) {
	for _, object := range blocks {
		drawRect(screen,
			float64(object.Body.X), float64(object.Body.Y), float64(object.Body.W), float64(object.Body.H),
			color.Gray{Y: object.gray})
	}
}

func drawBorders(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, width*widthFactor, float64(borderWidth), color.Gray{Y: 40})
	ebitenutil.DrawRect(screen, 0, float64(height-borderWidth), width*widthFactor, float64(height-borderWidth), color.Gray{Y: 40})
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
	text.Draw(screen, secs, arcadeFontLarge, width/2-len(secs)*50/2, height/2, color.Gray{Y: 200})
	text.Draw(screen, randomSeed.Code, arcadeFont, (width-len(randomSeed.Code)*10)/2, height/2+50, color.Gray{Y: 180})
}

func draw(screen *ebiten.Image) {
	drawBackground(screen)

	if !hud {
		drawGoal(screen, goal)
		drawPlayer(screen, player)
		drawBlocks(screen)
		drawBorders(screen)
		drawLevelCode(screen)
		drawTimer(screen)
	} else {
		drawHUD(screen)
	}

	debugInfo(screen)
	debugDrawPosition(screen)
}
