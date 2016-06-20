package main

import (
	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"

	"github.com/tanema/amore-examples/platformer/lib/gamera"
)

var (
	camera *gamera.Camera
	x      = float32(80)
	y      = float32(80)
)

func main() {
	amore.OnLoad = load
	amore.Start(update, draw)
}

func load() {
	camera = gamera.New(0, 0, 8000, 6000)
}

func update(dt float32) {
	if keyboard.IsDown(keyboard.KeyEscape) {
		amore.Quit()
	}
	if keyboard.IsDown(keyboard.KeyUp) {
		y -= 1
	}
	if keyboard.IsDown(keyboard.KeyDown) {
		y += 1
	}
	if keyboard.IsDown(keyboard.KeyLeft) {
		x -= 1
	}
	if keyboard.IsDown(keyboard.KeyRight) {
		x += 1
	}
	if keyboard.IsDown(keyboard.KeyS) {
		camera.Shake(1)
	}
	camera.SetPosition(x, y)
	camera.Update(dt)
}

func draw() {
	camera.Draw(func(l, t, w, h float32) {
		gfx.SetColor(0, 0, 255, 255)
		gfx.Circle(gfx.FILL, x, y, 20)

		gfx.SetColor(0, 255, 0, 255)
		gfx.Rect(gfx.FILL, 20, 20, 20, 20)
		gfx.Rect(gfx.FILL, 120, 20, 20, 20)
		gfx.Rect(gfx.FILL, 20, 120, 20, 20)
		gfx.Rect(gfx.FILL, 120, 120, 20, 20)
	})
}
