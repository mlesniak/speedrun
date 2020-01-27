package main

var currentScene *Scene

var scenes map[string]*Scene

type Scene struct {
	Init   func() error
	Reset  func() error
	Update func() error
	Draw   func() error
}

// We would define all scenes here?
// Create scene
// Add scene to list of scenes
// GameLoop uses Scenes, independent of content of a scene
// Allow to switch Scenes

// TODO Gamestate?
