package bump

import (
	"math"

	"github.com/tanema/amore/gfx"
)

const (
	DELTA = 1e-10 // floating-point margin of error
)

var (
	inf  = float32(math.Inf(1))
	ninf = float32(math.Inf(-1))
)

type (
	Entity interface {
		Send(event string, args ...interface{})
		Update(dt float32)
		Draw()
	}
	Body struct {
		Entity  Entity
		world   *World
		tag     string
		x       float32
		y       float32
		width   float32
		height  float32
		cells   []*Cell
		respMap map[string]string
	}
)

func (body *Body) move(x, y float32) (gx, gy float32, cols []*Collision) {
	actualX, actualY, collisions := body.check(x, y)
	body.update(actualX, actualY, body.width, body.height)
	return actualX, actualY, collisions
}

func (body *Body) check(goalX, goalY float32) (gx, gy float32, cols []*Collision) {
	collisions := []*Collision{}
	projected_cols := body.world.Project(body, goalX, goalY)
	visited := map[*Body]bool{body: true}

	for len(projected_cols) > 0 {
		collision := projected_cols[0]
		if _, ok := visited[collision.Other]; !ok {
			collisions = append(collisions, collision)
			response := body.world.responses[collision.RespType]
			goalX, goalY, projected_cols = response(collision, body, goalX, goalY)
			visited[collision.Other] = true
		}
	}

	return goalX, goalY, collisions
}

func (body *Body) update(x, y, w, h float32) {
	if body.x != x || body.y != y || body.width != w || body.height != h {
		for _, cell := range body.cells {
			cell.leave(body)
		}
		body.cells = []*Cell{}
		cl, ct, cw, ch := body.world.gridToCellRect(x, y, w, h)
		for cy := ct; cy <= ct+ch-1; cy++ {
			for cx := cl; cx <= cl+cw-1; cx++ {
				body.cells = append(body.cells, body.world.addToCell(body, cx, cy))
			}
		}
		body.x, body.y, body.width, body.height = x, y, w, h
	}
}

func (body *Body) remove() {
	for _, cell := range body.cells {
		cell.leave(body)
	}
}

func (body *Body) collide(other *Body, goalX, goalY float32) *Collision {
	if other == body {
		return nil
	}

	dx, dy := goalX-body.x, goalY-body.y
	diff := body.getDiff(other)

	var overlaps bool
	var ti, nx, ny float32

	if diff.containsPoint(0, 0) { // item was intersecting other
		px, py := diff.getNearestCorner(0, 0)
		// area of intersection
		wi := float32(math.Min(float64(body.width), math.Abs(float64(px))))
		hi := float32(math.Min(float64(body.height), math.Abs(float64(py))))
		ti = -wi * hi // ti is the negative area of intersection
		overlaps = true
	} else {
		ok, ti1, ti2, nx1, ny1, _, _ := diff.getSegmentIntersectionIndices(0, 0, dx, dy, ninf, inf)
		// item tunnels into other
		if ok && ti1 < 1 && (0 < ti1+DELTA || 0 == ti1 && ti2 > 0) {
			ti, nx, ny = ti1, nx1, ny1
			overlaps = false
		}
	}

	if ti == 0 {
		return nil
	}

	var tx, ty float32

	if overlaps {
		if dx == 0 && dy == 0 {
			// intersecting and not moving - use minimum displacement vector
			px, py := diff.getNearestCorner(0, 0)
			if math.Abs(float64(px)) < math.Abs(float64(py)) {
				py = 0
			} else {
				px = 0
			}
			nx, ny = sign(px), sign(py)
			tx, ty = body.x+px, body.y+py
		} else {
			// intersecting and moving - move in the opposite direction
			var ti1 float32
			var ok bool
			ok, ti1, _, nx, ny, _, _ = diff.getSegmentIntersectionIndices(0, 0, dx, dy, ninf, 1)
			if !ok {
				return nil
			}
			tx, ty = body.x+dx*ti1, body.y+dy*ti1
		}
	} else { // tunnel
		tx, ty = body.x+dx*ti, body.y+dy*ti
	}

	return &Collision{
		overlaps: overlaps,
		ti:       ti,
		Move:     Point{X: dx, Y: dy},
		Normal:   Point{X: nx, Y: ny},
		Touch:    Point{X: tx, Y: ty},
	}
}

// Calculates the minkowsky difference between 2 rects, which is another rect
func (body *Body) getDiff(other *Body) *Body {
	return &Body{
		x:      other.x - body.x - body.width,
		y:      other.y - body.y - body.height,
		width:  body.width + other.width,
		height: body.height + other.height,
	}
}

func (body *Body) containsPoint(px, py float32) bool {
	return px-body.x > DELTA && py-body.y > DELTA &&
		body.x+body.width-px > DELTA && body.y+body.height-py > DELTA
}

func (body *Body) getNearestCorner(px, py float32) (x, y float32) {
	return nearest(px, body.x, body.x+body.width), nearest(py, body.y, body.y+body.height)
}

func nearest(x, a, b float32) float32 {
	if math.Abs(float64(a-x)) < math.Abs(float64(b-x)) {
		return a
	}
	return b
}

func sign(x float32) float32 {
	if x > 0 {
		return 1
	}
	if x == 0 {
		return 0
	}
	return -1
}

// This is a generalized implementation of the liang-barsky algorithm, which also returns
// the normals of the sides where the segment intersects.
// Notice that normals are only guaranteed to be accurate when initially ti1, ti2 == -Inf, Inf
func (body *Body) getSegmentIntersectionIndices(x1, y1, x2, y2, ti1, ti2 float32) (bool, float32, float32, float32, float32, float32, float32) {
	dx, dy := x2-x1, y2-y1
	var nx, ny, p, q, nx1, ny1, nx2, ny2 float32
	for side := 1; side <= 4; side++ {
		if side == 1 { // left
			nx, ny, p, q = -1, 0, -dx, x1-body.x
		} else if side == 2 { // right
			nx, ny, p, q = 1, 0, dx, body.x+body.width-x1
		} else if side == 3 { // top
			nx, ny, p, q = 0, -1, -dy, y1-body.y
		} else { //bottom
			nx, ny, p, q = 0, 1, dy, body.y+body.height-y1
		}

		if p == 0 {
			if q <= 0 {
				return false, ti1, ti2, nx1, ny1, nx2, ny2
			}
		} else {
			r := q / p
			if p < 0 {
				if r > ti2 {
					return false, ti1, ti2, nx1, ny1, nx2, ny2
				} else if r > ti1 {
					ti1, nx1, ny1 = r, nx, ny
				}
			} else { // p > 0
				if r < ti1 {
					return false, ti1, ti2, nx1, ny1, nx2, ny2
				} else if r < ti2 {
					ti2, nx2, ny2 = r, nx, ny
				}
			}
		}
	}

	return true, ti1, ti2, nx1, ny1, nx2, ny2
}

func (body *Body) distanceTo(other *Body) float32 {
	dx := body.x - other.x + (body.width-other.width)/2
	dy := body.y - other.y + (body.height-other.height)/2
	return dx*dx + dy*dy
}

func (body *Body) DrawDebug() {
	gfx.SetColor(255, 0, 0, 100)
	gfx.Rect(gfx.FILL, body.x, body.y, body.width, body.height)
	gfx.SetColor(0, 255, 0, 100)
	gfx.Rect(gfx.LINE, body.x-1, body.y-1, body.width+2, body.height+2)
}

func (body *Body) Position() (x, y float32) {
	return body.x, body.y
}
