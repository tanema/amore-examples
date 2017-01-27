package main

import (
	"fmt"

	"github.com/tanema/amore-examples/asteroids/game"

	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/timer"
)

func main() {
	amore.OnLoad = load
	amore.Start(update, draw)
}

func load() {
	game.New()
}

func update(deltaTime float32) {
	if keyboard.IsDown(keyboard.KeyEscape) {
		amore.Quit()
	}
	game.Update(deltaTime)
}

func draw() {
	game.Draw()
	gfx.Print(fmt.Sprintf("fps: %v", timer.GetFPS()))
}
