package game

import (
	"github.com/tanema/amore/gfx"
)

var colors = map[string]*gfx.Color{
	"sky":    gfx.NewColor(114, 215, 238, 255),
	"tree":   gfx.NewColor(0, 81, 8, 255),
	"fog":    gfx.NewColor(0, 81, 8, 255),
	"road":   gfx.NewColor(107, 107, 107, 255),
	"grass":  gfx.NewColor(16, 170, 16, 255),
	"rumble": gfx.NewColor(85, 85, 85, 255),
	"lane":   gfx.NewColor(204, 204, 204, 255),
	"start":  gfx.NewColor(255, 255, 255, 255),
	"finish": gfx.NewColor(0, 0, 0, 255),
}
