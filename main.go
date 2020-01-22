package main

import (
	"fmt"
	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"log"
	"os"
)

type Object struct {
	gray uint8
	resolv.Rectangle
}

var player = Object{
	gray:      40,
	Rectangle: *resolv.NewRectangle(width/2, height/2, 20, 20),
}

func main() {
	if err := ebiten.Run(update, width, height, 1, title); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	checkExitKey()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	drawBackground(screen)
	debugInfo(screen)
	draw(screen, player)
	return nil
}

func checkExitKey() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
}

func debugInfo(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %d", ebiten.MaxTPS()))
}

func drawBackground(screen *ebiten.Image) {
	screen.Fill(color.Gray{Y: 100})
}

func draw(screen *ebiten.Image, object Object) {
	ebitenutil.DrawRect(screen, float64(object.X), float64(object.Y), float64(object.W), float64(object.H), color.Gray{Y: object.gray})
}
