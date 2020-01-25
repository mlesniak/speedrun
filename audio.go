package main

import (
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
	"github.com/markbates/pkger"
	"io/ioutil"
	"log"
)

var audioPlayer map[string]*audio.Player

func init() {
	audioPlayer = make(map[string]*audio.Player)
	audioContext, _ := audio.NewContext(44100)

	loadAudioAsset(audioContext, "countdown")
	loadAudioAsset(audioContext, "start")
	loadAudioAsset(audioContext, "goal")
	loadAudioAsset(audioContext, "highscore")
	loadAudioAsset(audioContext, "background")
}

// loadAudioAsset loads an ogg audio file with the fiven name from assets.
func loadAudioAsset(audioContext *audio.Context, name string) {
	// Load file.
	b, err := pkger.Open("/assets/" + name + ".ogg")
	mustAudio(err)
	defer b.Close()
	bs, err := ioutil.ReadAll(b)
	mustAudio(err)

	// Decode file.
	d, err := vorbis.Decode(audioContext, audio.BytesReadSeekCloser(bs))
	mustAudio(err)

	// Create player for future use.
	player, err := audio.NewPlayer(audioContext, d)
	mustAudio(err)
	audioPlayer[name] = player
}

// mustAudio checks if an error occured while loading a file.
func mustAudio(err error) {
	if err != nil {
		log.Fatal("Unable to load audio file:", err)
	}
}

// PlayAudio plays an audio file once.
func PlayAudio(name string) {
	PlayAudioTimes(name, 1)
}

// PlayAudioTimes plays an audio file multiple times.
func PlayAudioTimes(name string, times int) {
	player := audioPlayer[name]
	// If the audio is currently playing, rewind, i.e. do not play multiple streams.
	if player.IsPlaying() {
		player.Rewind()
	}

	go func() {
		// Reset after finish.
		defer player.Rewind()

		for i := 0; i < times; i++ {
			player.Play()
			for player.IsPlaying() {
				//time.Sleep(time.Millisecond * 10)
			}
			player.Rewind()
		}
	}()
}
