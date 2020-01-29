package main

import (
	"github.com/hajimehoshi/ebiten"
	"log"
)

var currentScene *Scene

// Scenes are mapped by identifiers.
var scenes = make(map[string]*Scene)

type Scene struct {
	Init   func()              // Called once while adding the bar.
	Reset  func()              // Called whenever the bar is reactivated.
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
