package main

import (
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
	"github.com/markbates/pkger"
	"io/ioutil"
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
	b, err := pkger.Open("/assets/" + name + ".ogg")

	if err != nil {
		panic(err)
	}
	defer b.Close()
	bs, _ := ioutil.ReadAll(b)
	d, _ := vorbis.Decode(audioContext, audio.BytesReadSeekCloser(bs))
	player, _ := audio.NewPlayer(audioContext, d)
	audioPlayer[name] = player
}

func playBackground(name string) {
	playBackgroundTimes(name, 1)
}

func playBackgroundTimes(name string, times int) {
	player := audioPlayer[name]
	if player.IsPlaying() {
		player.Rewind()
	}

	go func() {
		defer player.Rewind()
		for i := 0; i < times; i++ {
			player.Play()
			for player.IsPlaying() {
				// Wait...
			}
			player.Rewind()
		}
	}()
}
