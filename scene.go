package main

import (
	"github.com/hajimehoshi/ebiten"
	"log"
)

var currentScene *Scene

// Scenes are mapped by identifiers.
var scenes = make(map[string]*Scene)

// TODO Should Update() receive a pointer to a Gamestate?
// TODO What happens when a game scene is finished?
type Scene struct {
	Init   func()              // Called once while adding the scene.
	Reset  func()              // Called whenever the scene is reactivated.
	Update func()              // Called for state changes.
	Draw   func(*ebiten.Image) // Called for drawing state.
}

func AddScene(name string, scene *Scene) {
	scenes[name] = scene
	scene.Init()
}

func SetScene(name string) {
	scene, found := scenes[name]
	if !found {
		log.Fatal("Scene not found:", name)
	}
	currentScene = scene
	currentScene.Reset()
}

func GetCurrentScene() *Scene {
	return currentScene
}

// TODO Create two scenes: HUD and actual Game. Where should these be defined?
