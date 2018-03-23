package main

import (
	"github.com/tanema/amore"

	"github.com/tanema/gocraftmini/game"
)

func main() {
	world := game.NewWorld(149, 20, 300, false)
	amore.Start(world.Update, world.Draw)
}
