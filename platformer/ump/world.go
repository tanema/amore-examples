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
		bodies    []*Body
		responses map[string]Resp
	}
	Resp func(world *World, col *Collision, body *Body, goalX, goalY float32) (gx, gy float32, cols []*Collision)
)

func NewWorld(cellSize int) *World {
	world := &World{
		cellSize:  float32(cellSize),
		rows:      make(map[int]map[int]*Cell),
		responses: map[string]Resp{},
		bodies:    []*Body{},
	}

	world.AddResponse("touch", touchFilter)
	world.AddResponse("cross", crossFilter)
	world.AddResponse("slide", slideFilter)
	world.AddResponse("bounce", bounceFilter)

	return world
}

func (world *World) Add(tag string, left, top, w, h float32) *Body {
	body := &Body{
		world:   world,
		tag:     tag,
		w:       w,
		h:       h,
		cells:   []*Cell{},
		respMap: map[string]string{},
	}
	world.bodies = append(world.bodies, body)
	body.update(left, top, w, h)
	return body
}

func (world *World) Remove(body *Body) {
	body.remove()
	// TODO: delete from world.bodies
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

// TODO QuerySegment()

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
