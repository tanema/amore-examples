package main

import (
	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/mouse"

	"github.com/tanema/amore-examples/platformer/ump"
)

var (
	you   *ump.Body
	me    *ump.Body
	world *ump.World
)

func main() {
	world = ump.NewWorld(64)
	you = world.Add("block", 300, 300, 100, 100)
	me = world.Add("player", 0, 0, 32, 32)
	amore.Start(update, draw)
}

func update(deltaTime float32) {
	if keyboard.IsDown(keyboard.KeyEscape) {
		amore.Quit()
	}
	mx, my := mouse.GetPosition()
	me.Move(mx, my)
}

func draw() {
	world.DrawDebug(0, 0, 800, 600)

	gfx.SetColor(0, 0, 255, 255)
	x, y, w, h, _, _ := you.Extents()
	gfx.Rect(gfx.FILL, x, y, w, h)

	gfx.SetColorC(gfx.NewColor(0, 255, 0, 255))
	x, y, w, h, _, _ = me.Extents()
	gfx.Rect(gfx.FILL, x, y, w, h)
}
