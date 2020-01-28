package main

// Pack all assets into pkged.go for single file distribution.
//
//go:generate go get github.com/markbates/pkger/cmd/pkger
//go:generate pkger -include /assets

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mlesniak/speedrun/internal/seed"
	"log"
	"os"
)

func main() {
	// initalizeGame()

	randomSeed = seed.New()
	if len(os.Args) > 1 {
		randomSeed = seed.NewPreset(os.Args[1])
	}

	AddScene("hud", hudScene)
	AddScene("game", gameScene)
	SetScene("game")

	// Start game loop.
	if err := ebiten.Run(GameLoop, width, height, 1.0, title); err != nil {
		log.Fatal("Unable to start game: ", err)
	}
}

// GameLoop is the main game loop, updating the current game updateState and (optionally) drawing it.
func GameLoop(screen *ebiten.Image) error {
	//background.Update()
	//
	//// Legacy objects.
	//updateState()
	//
	//if ebiten.IsDrawingSkipped() {
	//	return nil
	//}
	//
	//background.Draw(screen)
	//// Legacy objects.
	//drawState(screen)
	//return nil

	GetCurrentScene().Update()
	if !ebiten.IsDrawingSkipped() {
		GetCurrentScene().Draw(screen)
	}

	return nil
}
