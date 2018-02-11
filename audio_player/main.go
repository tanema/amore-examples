package main

import (
	"fmt"
	"time"

	"github.com/tanema/amore"
	"github.com/tanema/amore/audio"
	"github.com/tanema/amore/gfx"

	"github.com/tanema/amore-examples/audio_player/ui"
)

var (
	elements     = []ui.Element{}
	playing      = false
	currentTrack = 0
	trackNames   = []string{"audio/test.wav", "audio/test.mp3", "audio/test.ogg", "audio/test.flac"}
	tracks       = []*audio.Source{}
)

const buttonSize float32 = 64

func main() {
	elements = append(elements,
		ui.NewButton(0*buttonSize, 0, buttonSize, buttonSize, ui.BackImg, ui.Clear, func(button *ui.Button) {
			tracks[currentTrack].Stop()
			currentTrack = (currentTrack - 1 + len(tracks)) % len(tracks)
			tracks[currentTrack].Play()
		}),
		ui.NewButton(1*buttonSize, 0, buttonSize, buttonSize, ui.RewindImg, ui.Clear, func(button *ui.Button) {
			tracks[currentTrack].Seek(tracks[currentTrack].Tell() - (5 * time.Second))
		}),
		ui.NewButton(2*buttonSize, 0, buttonSize, buttonSize, ui.PlayImg, ui.Clear, func(button *ui.Button) {
			if tracks[currentTrack].IsPlaying() {
				tracks[currentTrack].Pause()
			} else {
				tracks[currentTrack].Play()
			}

			if tracks[currentTrack].IsPlaying() {
				button.SetImage(ui.PauseImg)
			} else {
				button.SetImage(ui.PlayImg)
			}
		}),
		ui.NewButton(3*buttonSize, 0, buttonSize, buttonSize, ui.StopImg, ui.Clear, func(button *ui.Button) {
			tracks[currentTrack].Stop()
		}),
		ui.NewButton(4*buttonSize, 0, buttonSize, buttonSize, ui.ForwardImg, ui.Clear, func(button *ui.Button) {
			tracks[currentTrack].Seek(tracks[currentTrack].Tell() + (5 * time.Second))
		}),
		ui.NewButton(5*buttonSize, 0, buttonSize, buttonSize, ui.NextImg, ui.Clear, func(button *ui.Button) {
			tracks[currentTrack].Stop()
			currentTrack = (currentTrack + 1) % len(tracks)
			tracks[currentTrack].Play()
		}),
	)
	amore.OnLoad = load
	amore.Start(update, draw)
}

func load() {
	for _, name := range trackNames {
		source, err := audio.NewSource(name, false)
		if err != nil {
			panic(err)
		}
		tracks = append(tracks, source)
	}
}

func update(deltaTime float32) {
	for _, element := range elements {
		element.Update(deltaTime)
	}
}

func draw() {
	gfx.SetBackgroundColor(255, 255, 255, 255)
	for _, element := range elements {
		element.Draw()
	}
	gfx.SetColorC(ui.Black)
	gfx.Print(fmt.Sprintf("%v", tracks[currentTrack].GetDuration()), 0, 64)
	gfx.Print(fmt.Sprintf("%v", tracks[currentTrack].Tell()), 0, 128)
	gfx.Print(trackNames[currentTrack], 200, 128)
}
