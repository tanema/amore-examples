package ump

import (
	"math"
)

type Grid struct {
	cellSize float32
	rows     map[int]map[int]*Cell
}

func newGrid(cellSize int) *Grid {
	return &Grid{
		cellSize: float32(cellSize),
		rows:     make(map[int]map[int]*Cell),
	}
}

func (grid *Grid) update(body *Body) {
	for _, cell := range body.cells {
		cell.leave(body)
	}
	body.cells = []*Cell{}
	cl, ct, cw, ch := grid.toCellRect(body.x, body.y, body.w, body.h)
	for cy := ct; cy <= ct+ch-1; cy++ {
		for cx := cl; cx <= cl+cw-1; cx++ {
			grid.cellAt(float32(cx), float32(cy), true).enter(body)
		}
	}
}

func (grid *Grid) cellsInRect(l, t, w, h float32) []*Cell {
	cl, ct, cw, ch := grid.toCellRect(l, t, w, h)
	cells := []*Cell{}
	for cy := ct; cy <= ct+ch-1; cy++ {
		row, ok := grid.rows[cy]
		if ok {
			for cx := cl; cx <= cl+cw-1; cx++ {
				cell, ok := row[cx]
				if ok {
					cells = append(cells, cell)
				}
			}
		}
	}
	return cells
}

func (grid *Grid) toCellRect(x, y, w, h float32) (cx, cy, cw, ch int) {
	cx, cy = grid.cellCoordsAt(x, y)
	cr, cb := int(math.Ceil(float64((x+w)/grid.cellSize))), int(math.Ceil(float64((y+h)/grid.cellSize)))
	return cx, cy, cr - cx, cb - cy
}

func (grid *Grid) cellCoordsAt(x, y float32) (cx, cy int) {
	return int(math.Floor(float64(x / grid.cellSize))), int(math.Floor(float64(y / grid.cellSize)))
}

func (grid *Grid) cellAt(x, y float32, cellCoords bool) *Cell {
	var cx, cy int
	if cellCoords == true {
		cx, cy = int(x), int(y)
	} else {
		cx, cy = grid.cellCoordsAt(x, y)
	}
	row, ok := grid.rows[cy]
	if !ok {
		grid.rows[cy] = make(map[int]*Cell)
		row = grid.rows[cy]
	}
	cell, ok := row[cx]
	if !ok {
		row[cx] = &Cell{bodies: make(map[uint32]*Body)}
		cell = row[cx]
	}
	return cell
}

func (grid *Grid) getCellsTouchedBySegment(x1, y1, x2, y2 float32) []*Cell {
	cells := []*Cell{}
	visited := map[*Cell]bool{}

	grid.traceRay(x1, y1, x2, y2, func(cx, cy int) {
		cell := grid.cellAt(float32(cx), float32(cy), true)
		if _, found := visited[cell]; found {
			return
		}
		visited[cell] = true
		cells = append(cells, cell)
	})

	return cells
}

// traceRay* functions are based on "A Fast Voxel Traversal Algorithm for Ray Tracing",
// by John Amanides and Andrew Woo - http://www.cse.yorku.ca/~amana/research/grid.pdf
// It has been modified to include both cells when the ray "touches a grid corner",
// and with a different exit condition
func (grid *Grid) rayStep(ct, t1, t2 float32) (int, float32, float32) {
	v := t2 - t1
	delta := grid.cellSize / v
	if v > 0 {
		return 1, delta, delta * (1.0 - frac(t1/grid.cellSize))
	} else if v < 0 {
		return -1, -delta, -delta * frac(t1/grid.cellSize)
	} else {
		return 0, inf, inf
	}
}

func (grid *Grid) traceRay(x1, y1, x2, y2 float32, f func(cx, cy int)) {
	cx1, cy1 := grid.cellCoordsAt(x1, y1)
	cx2, cy2 := grid.cellCoordsAt(x2, y2)
	stepX, dx, tx := grid.rayStep(float32(cx1), x1, x2)
	stepY, dy, ty := grid.rayStep(float32(cy1), y1, y2)
	cx, cy := cx1, cy1

	f(cx, cy)

	// The default implementation had an infinite loop problem when
	// approaching the last cell in some occassions. We finish iterating
	// when we are *next* to the last cell
	for abs(float32(cx-cx2))+abs(float32(cy-cy2)) > 1 {
		if tx < ty {
			tx += dx
			cx += stepX
		} else {
			// Addition: include both cells when going through corners
			if tx == ty {
				f(cx+stepX, cy)
			}
			ty += dy
			cy += stepY
		}
		f(cx, cy)
	}

	// If we have not arrived to the last cell, use it
	if cx != cx2 || cy != cy2 {
		f(cx2, cy2)
	}
}
