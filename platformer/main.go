package main

import (
	"github.com/tanema/amore-examples/platformer/game"

	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"
)

func main() {
	amore.OnLoad = game.New
	amore.Start(update, draw)
}

func update(deltaTime float32) {
	if keyboard.IsDown(keyboard.KeyEscape) {
		amore.Quit()
	}
	game.Update(deltaTime)
}

func draw() {
	game.Draw()
}
