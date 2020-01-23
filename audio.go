package main

import (
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var audioPlayer map[string]*audio.Player
var context *audio.Context

func init() {
	audioPlayer = make(map[string]*audio.Player)

	context, _ = audio.NewContext(44100)
	loadAudio(context, "countdown")
	loadAudio(context, "start")
	loadAudio(context, "goal")
	loadAudio(context, "highscore")
	loadAudio(context, "background")
}

func loadAudio(audioContext *audio.Context, name string) {
	file, _ := ebitenutil.OpenFile("assets/" + name + ".ogg")
	d, _ := vorbis.Decode(audioContext, file)
	player, _ := audio.NewPlayer(audioContext, d)
	audioPlayer[name] = player
}

func playBackground(name string) {
	player := audioPlayer[name]
	go func() {
		player.Play()
		defer player.Rewind()
	}()
}
