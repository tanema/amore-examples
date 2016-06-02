package game

import (
	"github.com/tanema/amore-examples/platformer/lib/gamera"

	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"
)

const (
	updateRadius float32 = 100 // how "far away from the camera" things stop being updated
	instructions         = `
  bump.lua demo
    left,right: move
    up:     jump/fly
    return: reset map
    tab:    toggle debug info`
)

var (
	width     float32 = 4000
	height    float32 = 2000
	camera    *gamera.Camera
	game_map  *Map
	drawDebug = false // draw bump's debug info, fps and memory
)

func New() {
	keyboard.OnKeyUp = keyup
	camera = gamera.New(0, 0, width, height)
	game_map = newMap(width, height, camera)
}

func Update(dt float32) {
	l, t, w, h := camera.GetVisible()
	l, t, w, h = l-updateRadius, t-updateRadius, w+updateRadius*2, h+updateRadius*2
	game_map.update(dt, l, t, w, h)
	camera.SetPosition(game_map.player.getCenter())
	camera.Update(dt)
}

func Draw() {
	camera.Draw(func(l, t, w, h float32) {
		game_map.draw(drawDebug, l, t, w, h)
	})
	gfx.SetColor(255, 255, 255, 255)
	gfx.Print(instructions, gfx.GetWidth()-200, 10)
}

func keyup(key keyboard.Key) {
	if key == keyboard.KeyTab {
		drawDebug = !drawDebug
	}
	if key == keyboard.KeyReturn {
		game_map.reset()
	}
}
