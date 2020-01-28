package main

import "log"

var currentScene *Scene

// Scenes are mapped by identifiers.
var scenes = make(map[string]*Scene)

// TODO Should Update() receive a pointer to a Gamestate?
type Scene struct {
	Init   func() error // Called once while adding the scene.
	Reset  func() error // Called whenever the scene is reactivated.
	Update func() error // Called for state changes.
	Draw   func() error // Called for drawing state.
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
