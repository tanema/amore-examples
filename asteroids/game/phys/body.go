package phys

import (
	"github.com/tanema/amore/gfx"
)

type Body struct {
	world      *World
	Name       string
	inital     []float32
	points     []float32
	cells      []*Cell
	Collidable Collidable
}

type Collidable interface {
	Destroy(force bool)
}

func newBody(world *World, collidable Collidable, name string, points []float32) *Body {
	return &Body{
		Name:       name,
		world:      world,
		Collidable: collidable,
		inital:     points,
		points:     make([]float32, len(points)),
	}
}

func (body *Body) Move(x, y, rot, scale float32) []*Body {
	if body.Collidable != nil {
		for _, cell := range body.cells {
			cell.leave(body)
		}
	}

	mat := &Matrix{}
	mat.translate(x, y, rot, scale)
	for i := 0; i < len(body.inital); i += 2 {
		body.points[i], body.points[i+1] = mat.multiply(body.inital[i], body.inital[i+1])
	}

	if body.Collidable != nil {
		body.cells = []*Cell{}
		others := []*Body{}
		others_map := map[*Body]bool{}
		for i := 0; i < len(body.points); i += 2 {
			cell := body.world.CellAt(body.points[i], body.points[i+1])

			for _, other := range cell.bodies {
				if other != body && other.pointInside(body.points[i], body.points[i+1]) {
					if _, ok := others_map[other]; !ok {
						others = append(others, other)
					}
					others_map[other] = true
				}
			}

			body.cells = append(body.cells, cell)
			cell.enter(body)
		}

		return others
	}

	return []*Body{}
}

func (body *Body) pointInside(x, y float32) bool {
	j := 2
	oddNodes := false
	for i := 0; i < len(body.points); i += 2 {
		y0 := body.points[i+1]
		y1 := body.points[j+1]
		if (y0 < y && y1 >= y) || (y1 < y && y0 >= y) {
			if body.points[i]+(y-y0)/(y1-y0)*(body.points[j]-body.points[i]) < x {
				oddNodes = !oddNodes
			}
		}
		j += 2
		if j == len(body.points) {
			j = 0
		}
	}
	return oddNodes
}

func (body *Body) GetPoints() []float32 {
	return body.points
}

func (body *Body) Remove() {
	for _, cell := range body.cells {
		cell.leave(body)
	}
}

func (body *Body) Draw() {
	gfx.SetColor(255, 0, 0, 100)
	for _, cell := range body.cells {
		gfx.Rect(gfx.LINE, cell.x, cell.y, cell.width, cell.height)
	}
	gfx.SetColor(255, 255, 255, 255)
}
