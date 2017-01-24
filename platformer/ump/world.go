package ump

import (
	"math"
	"sort"

	"github.com/tanema/amore/gfx"
)

const defaultFilter = "slide"

type (
	World struct {
		grid      *Grid
		responses map[string]Resp
	}
	Resp func(world *World, col *Collision, body *Body, goalX, goalY float32) (gx, gy float32, cols []*Collision)
)

func NewWorld(cellSize int) *World {
	world := &World{
		grid:      newGrid(cellSize),
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

func (world *World) QueryRect(x, y, w, h float32, tags ...string) []*Body {
	return world.getBodiesInCells(world.grid.cellsInRect(x, y, w, h), tags...)
}

func (world *World) QueryPoint(x, y float32, tags ...string) []*Body {
	bodies := []*Body{}
	cell := world.grid.cellAt(x, y, false)
	if cell == nil {
		return []*Body{}
	}
	for _, body := range cell.bodies {
		if body.HasTag(tags...) && body.containsPoint(x, y) {
			bodies = append(bodies, body)
		}
	}
	return bodies
}

func (world *World) QuerySegment(x1, y1, x2, y2 float32, tags ...string) []*Body {
	bodies := []*Body{}
	visited := map[*Body]bool{}
	cells := world.grid.getCellsTouchedBySegment(x1, y1, x2, y2)
	bodiesOnSegment := world.getBodiesInCells(cells)
	distances := map[uint32]float32{}
	for _, body := range bodiesOnSegment {
		if _, ok := visited[body]; !ok && body.HasTag(tags...) {
			visited[body] = true
			fraction, _, _ := body.getRayIntersectionFraction(x1, y1, x2-x1, y2-y1)
			if fraction != inf {
				bodies = append(bodies, body)
				distances[body.ID] = fraction
			}
		}
	}

	By(func(b1, b2 *Body) bool {
		return distances[b1.ID] < distances[b2.ID]
	}).Sort(bodies)

	return bodies
}

func (world *World) getBodiesInCells(cells []*Cell, tags ...string) []*Body {
	dict := make(map[uint32]bool)
	bodies := []*Body{}
	for _, cell := range cells {
		for id, body := range cell.bodies {
			if _, ok := dict[id]; !ok && body.HasTag(tags...) {
				bodies = append(bodies, body)
				dict[id] = true
			}
		}
	}
	return bodies
}

func (world *World) Project(body *Body, goalX, goalY float32) []*Collision {
	collisions := []*Collision{}

	tl := float32(math.Min(float64(goalX), float64(body.x)))
	tt := float32(math.Min(float64(goalY), float64(body.y)))
	tr := float32(math.Max(float64(goalX+body.w), float64(body.x+body.w)))
	tb := float32(math.Max(float64(goalY+body.h), float64(body.y+body.h)))

	visited := map[*Body]bool{}
	bodies := world.getBodiesInCells(world.grid.cellsInRect(tl, tt, tr-tl, tb-tt))
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
	cellSize := world.grid.cellSize
	cl, ct, cw, ch := world.grid.toCellRect(l, t, w, h)
	for cy := ct; cy <= ct+ch-1; cy++ {
		row, ok := world.grid.rows[cy]
		if ok {
			for cx := cl; cx <= cl+cw-1; cx++ {
				cell, ok := row[cx]
				if ok {
					l, t, w, h := float32(cx)*cellSize, float32(cy)*cellSize, cellSize, cellSize
					intensity := cell.itemCount*12 + 16
					gfx.SetColor(255, 255, 255, float32(intensity))
					gfx.Rect(gfx.FILL, l, t, w, h)
					gfx.SetColor(255, 255, 255, 10)
					gfx.Rect(gfx.LINE, l, t, w, h)
				}
			}
		}
	}
}
