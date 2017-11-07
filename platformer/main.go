package main

import (
	"fmt"
	"runtime"

	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/timer"
	"github.com/tanema/lense"

	"github.com/tanema/amore-examples/platformer/game"
)

var (
	width   float32 = 4000
	height  float32 = 2000
	camera  *lense.Camera
	gameMap *game.Map
)

func main() {
	amore.OnLoad = onLoad
	amore.Start(update, draw)
}

func onLoad() {
	keyboard.OnKeyUp = keypress
	camera = lense.New()
	gameMap = game.NewMap(width, height, camera)
}

func update(dt float32) {
	l, t, w, h := camera.GetVisible()
	gameMap.Update(dt, l, t, w, h)
	camera.LookAt(gameMap.Player.GetCenter())
	camera.Update(dt)
}

func draw() {
	camera.Draw(gameMap.Draw)
	gfx.SetColor(255, 255, 255, 255)
	w, h := gfx.GetWidth(), gfx.GetHeight()
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	gfx.Printf(fmt.Sprintf("fps: %v, mem: %vKB", timer.GetFPS(), stats.HeapAlloc/1000000), 200, gfx.AlignRight, w-200, h-40)
}

func keypress(key keyboard.Key) {
	switch key {
	case keyboard.KeyEscape:
		amore.Quit()
	case keyboard.KeyTab:
		gameMap.ToggleDebug()
	case keyboard.KeyReturn:
		gameMap.Reset()
	}
}
