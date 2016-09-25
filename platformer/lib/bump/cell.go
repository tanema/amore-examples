package bump

import (
	"github.com/tanema/amore/gfx"
)

type Cell struct {
	world  *World
	cx     int
	cy     int
	bodies []*Body
}

func (cell *Cell) enter(body *Body) {
	cell.leave(body) // make sure this body is uniqe
	cell.bodies = append(cell.bodies, body)
}

func (cell *Cell) leave(body *Body) {
	for i, b := range cell.bodies {
		if body == b {
			cell.bodies = append(cell.bodies[:i], cell.bodies[i+1:]...)
			return
		}
	}
}

func (cell *Cell) DrawDebug(cellSize float32) {
	gfx.SetLineWidth(2)
	if len(cell.bodies) > 0 {
		gfx.SetColor(0, 255, 0, 200)
	} else {
		gfx.SetColor(0, 0, 255, 200)
	}
	gfx.Rect(gfx.LINE, float32(cell.cx)*cellSize, float32(cell.cy)*cellSize, cellSize, cellSize)
}
