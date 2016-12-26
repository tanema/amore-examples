package ump

import (
	"math"
	"sort"
)

const defaultFilter = "slide"

type (
	World struct {
		cellSize  float32
		rows      map[int]map[int]*Cell
		responses map[string]Resp
	}
	Resp func(world *World, col *Collision, body *Body, goalX, goalY float32) (gx, gy float32, cols []*Collision)
)

func NewWorld(cellSize int) *World {
	world := &World{
		cellSize:  float32(cellSize),
		rows:      make(map[int]map[int]*Cell),
		responses: map[string]Resp{},
	}

	world.AddResponse("touch", touchFilter)
	world.AddResponse("cross", crossFilter)
	world.AddResponse("slide", slideFilter)
	world.AddResponse("bounce", bounceFilter)

	return world
}

func (world *World) Add(tag string, left, top, w, h float32) *Body {
	return newBody(world, tag, left, top, w, h)
}

func (world *World) addToCell(body *Body, cx, cy int) *Cell {
	row, ok := world.rows[cy]
	if !ok {
		world.rows[cy] = make(map[int]*Cell)
		row = world.rows[cy]
	}
	cell, ok := row[cx]
	if !ok {
		row[cx] = &Cell{world: world, cx: cx, cy: cy}
		cell = row[cx]
	}
	cell.enter(body)
	return cell
}

func (world *World) QueryRect(x, y, w, h float32) []*Body {
	return world.getBodiesInCells(world.cellsInRect(x, y, w, h))
}

func (world *World) QueryPoint(x, y float32) []*Body {
	bodies := []*Body{}
	cx, cy := world.cellAt(x, y)
	var cell *Cell

	row, ok := world.rows[cy]
	if ok {
		cell, _ = row[cx]
	}
	if cell == nil {
		return []*Body{}
	}

	for _, body := range cell.bodies {
		if body.containsPoint(x, y) {
			bodies = append(bodies, body)
		}
	}
	return bodies
}

func (world *World) QuerySegment(x1, y1, x2, y2 float32) []*Body {
	bodies := []*Body{}
	visited := map[*Body]bool{}
	for _, body := range world.getBodiesInCells(world.getCellsTouchedBySegment(x1, y1, x2, y2)) {
		if _, ok := visited[body]; ok {
			visited[body] = true
			fraction, _, _ := body.getRayIntersectionFraction(x1, y1, x2, y2)
			if fraction != inf {
				bodies = append(bodies, body)
			}
		}
	}
	return bodies
}

func (world *World) getCellsTouchedBySegment(x1, y1, x2, y2 float32) []*Cell {
	cells := []*Cell{}
	visited := map[*Cell]bool{}

	world.traceRay(x1, y1, x2, y2, func(cx, cy int) {
		row, ok := world.rows[cy]
		if !ok {
			return
		}
		cell, ok := row[cx]
		if _, found := visited[cell]; found || !ok {
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
func (world *World) rayStep(ct, t1, t2 float32) (step int, dx, dy float32) {
	v := t2 - t1
	if v > 0 {
		return 1, world.cellSize / v, ((ct+v)*world.cellSize - t1) / v
	} else if v < 0 {
		return -1, -world.cellSize / v, ((ct+v-1)*world.cellSize - t1) / v
	} else {
		return 0, inf, inf
	}
}

func (world *World) traceRay(x1, y1, x2, y2 float32, f func(cx, cy int)) {
	cx1, cy1 := world.cellAt(x1, y1)
	cx2, cy2 := world.cellAt(x2, y2)
	stepX, dx, tx := world.rayStep(float32(cx1), x1, x2)
	stepY, dy, ty := world.rayStep(float32(cy1), y1, y2)
	cx, cy := cx1, cy1

	f(cx, cy)

	// The default implementation had an infinite loop problem when
	// approaching the last cell in some occassions. We finish iterating
	// when we are *next* to the last cell
	for math.Abs(float64(cx-cx2))+math.Abs(float64(cy-cy2)) > 1 {
		if tx < ty {
			tx, cx = tx+dx, cx+stepX
			f(cx, cy)
		} else {
			// Addition: include both cells when going through corners
			if tx == ty {
				f(cx+stepX, cy)
			}
			ty, cy = ty+dy, cy+stepY
			f(cx, cy)
		}
	}

	// If we have not arrived to the last cell, use it
	if cx != cx2 || cy != cy2 {
		f(cx2, cy2)
	}
}

func (world *World) cellsInRect(l, t, w, h float32) []*Cell {
	cl, ct, cw, ch := world.gridToCellRect(l, t, w, h)
	cells := []*Cell{}
	for cy := ct; cy <= ct+ch-1; cy++ {
		row, ok := world.rows[cy]
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

func (world *World) getBodiesInCells(cells []*Cell) []*Body {
	dict := make(map[*Body]bool)
	bodies := []*Body{}
	for _, cell := range cells {
		for _, body := range cell.bodies {
			if _, ok := dict[body]; !ok {
				bodies = append(bodies, body)
				dict[body] = true
			}
		}
	}
	return bodies
}

func (world *World) gridToCellRect(x, y, w, h float32) (cx, cy, cw, ch int) {
	cx, cy = world.cellAt(x, y)
	cr, cb := int(math.Ceil(float64((x+w)/world.cellSize))), int(math.Ceil(float64((y+h)/world.cellSize)))
	return cx, cy, cr - cx, cb - cy
}

func (world *World) cellAt(x, y float32) (cx, cy int) {
	return int(math.Floor(float64(x / world.cellSize))), int(math.Floor(float64(y / world.cellSize)))
}

func (world *World) Project(body *Body, goalX, goalY float32) []*Collision {
	collisions := []*Collision{}

	tl := float32(math.Min(float64(goalX), float64(body.x)))
	tt := float32(math.Min(float64(goalY), float64(body.y)))
	tr := float32(math.Max(float64(goalX+body.w), float64(body.x+body.w)))
	tb := float32(math.Max(float64(goalY+body.h), float64(body.y+body.h)))

	visited := map[*Body]bool{}
	bodies := world.getBodiesInCells(world.cellsInRect(tl, tt, tr-tl, tb-tt))
	for _, other := range bodies {
		if _, ok := visited[other]; !ok {
			visited[other] = true
			if col := body.collide(other, goalX, goalY); col != nil {
				collisions = append(collisions, col)
			}
		}
	}

	sort.Sort(CollisionsByDistance(collisions))

	return collisions
}

func (world *World) AddResponse(name string, response Resp) {
	world.responses[name] = response
}

func (world *World) DrawDebug(l, t, w, h float32) {
	for _, cell := range world.cellsInRect(l, t, w, h) {
		cell.DrawDebug(world.cellSize)
	}
}
