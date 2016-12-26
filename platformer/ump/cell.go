package ump

import (
	"fmt"

	"github.com/tanema/amore/gfx"
)

type Cell struct {
	world  *World
	cx     int
	cy     int
	bodies []*Body
	text   *gfx.Text
}

func (cell *Cell) enter(body *Body) {
	cell.leave(body) // make sure this body is uniqe
	cell.bodies = append(cell.bodies, body)
	cell.text, _ = gfx.NewText(gfx.GetFont(), fmt.Sprintf("%v", len(cell.bodies)))
}

func (cell *Cell) leave(body *Body) {
	for i, b := range cell.bodies {
		if body == b {
			cell.bodies = append(cell.bodies[:i], cell.bodies[i+1:]...)
			cell.text, _ = gfx.NewText(gfx.GetFont(), fmt.Sprintf("%v", len(cell.bodies)))
			return
		}
	}
}

func (cell *Cell) DrawDebug(cellSize float32) {
	gfx.SetLineWidth(2)
	x, y := float32(cell.cx)*cellSize, float32(cell.cy)*cellSize
	half := (cellSize / 2) - 5
	gfx.SetColor(100, 100, 100, 175)
	cell.text.Draw(x+half, y+half)
	gfx.Rect(gfx.LINE, x, y, cellSize, cellSize)
}
