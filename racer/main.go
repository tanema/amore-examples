package main

import (
	"github.com/tanema/amore"
	"github.com/tanema/amore/keyboard"

	"github.com/tanema/amore-examples/racer/game"
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
