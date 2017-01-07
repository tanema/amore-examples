package game

import "github.com/tanema/amore/gfx"

func drawFilledRectangle(l, t, w, h, r, g, b float32) {
	gfx.SetColor(r, g, b, 100)
	gfx.Rect(gfx.FILL, l, t, w, h)
	gfx.SetColor(r, g, b, 255)
	gfx.Rect(gfx.LINE, l, t, w, h)
}
