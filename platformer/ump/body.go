package ump

import (
	"github.com/tanema/amore/gfx"
)

type (
	Body struct {
		world   *World
		tag     string
		x       float32
		y       float32
		w       float32
		h       float32
		cells   []*Cell
		respMap map[string]string // TODO: Add api for response map
	}
)

func (body *Body) Move(x, y float32) (gx, gy float32, cols []*Collision) {
	actualX, actualY, collisions := body.check(x, y)
	body.update(actualX, actualY, body.w, body.h)
	return actualX, actualY, collisions
}

func (body *Body) check(goalX, goalY float32) (gx, gy float32, cols []*Collision) {
	collisions := []*Collision{}
	projected_cols := body.world.Project(body, goalX, goalY)
	visited := map[*Body]bool{body: true}

	for len(projected_cols) > 0 {
		collision := projected_cols[0]
		if _, ok := visited[collision.Body]; !ok {
			collisions = append(collisions, collision)
			response := body.world.responses[collision.RespType]
			goalX, goalY, projected_cols = response(body.world, collision, body, goalX, goalY)
			visited[collision.Body] = true
		} else {
			projected_cols = projected_cols[1:]
		}
	}

	return goalX, goalY, collisions
}

func (body *Body) update(x, y, w, h float32) {
	if body.x != x || body.y != y || body.w != w || body.h != h {
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
		body.x, body.y, body.w, body.h = x, y, w, h
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
	respType, ok := body.respMap[other.tag]
	if !ok {
		respType = defaultFilter
	}

	collision := &Collision{
		Body:     other,
		RespType: respType,
		Distance: body.distanceTo(other),
		Move:     Point{X: dx, Y: dy},
	}

	// intersecting and not moving - use minimum displacement vector
	if diff.containsPoint(0, 0) && dx == 0 && dy == 0 {
		px, py := diff.getNearestCorner(0, 0)
		collision.Intersection = -min(body.w, abs(px)) * min(body.h, abs(py))
		if abs(px) < abs(py) {
			py = 0
		} else {
			px = 0
		}
		collision.Normal = Point{X: sign(px), Y: sign(py)}
		collision.Touch = Point{
			X: body.x + px + (collision.Normal.X * delta),
			Y: body.y + py + (collision.Normal.Y * delta),
		}
	} else {
		collision.Intersection, collision.Normal.X, collision.Normal.Y = diff.getRayIntersectionFraction(0, 0, dx, dy)
		if collision.Intersection == inf { //no intersection, no collision
			return nil
		}
		collision.Touch = Point{
			X: body.x + dx*collision.Intersection + (collision.Normal.X * delta),
			Y: body.y + dy*collision.Intersection + (collision.Normal.Y * delta),
		}
	}

	return collision
}

// Calculates the minkowski difference between 2 rects, which is another rect
func (body *Body) getDiff(other *Body) *Body {
	return &Body{
		x: other.x - body.x - body.w,
		y: other.y - body.y - body.h,
		w: body.w + other.w,
		h: body.h + other.h,
	}
}

func (body *Body) containsPoint(px, py float32) bool {
	return body.x < px && body.x+body.w > px &&
		body.y < py && body.y+body.h > py
}

func (body *Body) getNearestCorner(px, py float32) (x, y float32) {
	return nearest(px, body.x, body.x+body.w), nearest(py, body.y, body.y+body.h)
}

func (body *Body) getRayIntersectionFraction(ox, oy, dx, dy float32) (fraction, nx, ny float32) {
	vec := []float32{ox, oy, ox + dx, oy + dy}
	fraction = inf
	right := body.x + body.w
	bottom := body.y + body.h

	rayTests := [4][6]float32{
		{-1, 0, body.x, body.y, body.x, bottom},
		{0, 1, body.x, bottom, right, bottom},
		{1, 0, right, bottom, right, body.y},
		{0, -1, right, body.y, body.x, body.y},
	}

	// for each of the AABB's four edges calculate the minimum fraction of "direction"
	// in order to find where the ray FIRST intersects the AABB (if it ever does)
	for _, side := range rayTests {
		x := getRayIntersectionFractionOfFirstRay(vec, side[2:])
		if x < fraction {
			fraction = x
			nx, ny = side[0], side[1]
		}
	}

	return fraction, nx, ny
}

// taken from https://github.com/pgkelley4/line-segments-intersect/blob/master/js/line-segments-intersect.js
// which was adapted from http://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect
// returns the point where they intersect (if they intersect)
// returns inf if they don't intersect
func getRayIntersectionFractionOfFirstRay(vec1, vec2 []float32) float32 {
	rx, ry := vec1[2]-vec1[0], vec1[3]-vec1[1]
	sx, sy := vec2[2]-vec2[0], vec2[3]-vec2[1]

	numerator := crossProduct(vec2[0]-vec1[0], vec2[1]-vec1[1], rx, ry)
	denominator := crossProduct(rx, ry, sx, sy)

	// lines are parallel or the lines are co-linear
	if denominator == 0 {
		return inf
	}

	u := numerator / denominator
	t := crossProduct(vec2[0]-vec1[0], vec2[1]-vec1[1], sx, sy) / denominator
	if (t >= 0) && (t <= 1) && (u >= 0) && (u <= 1) {
		return t
	}

	return inf
}

func (body *Body) distanceTo(other *Body) float32 {
	dx := body.x - other.x + (body.w-other.w)/2
	dy := body.y - other.y + (body.h-other.h)/2
	return dx*dx + dy*dy
}

func (body *Body) DrawDebug() {
	gfx.SetColor(255, 0, 0, 100)
	gfx.Rect(gfx.FILL, body.x, body.y, body.w, body.h)
	gfx.SetColor(0, 255, 0, 100)
	gfx.Rect(gfx.LINE, body.x-1, body.y-1, body.w+2, body.h+2)
}

func (body *Body) Position() (x, y float32) {
	return body.x, body.y
}

func (body *Body) Extents() (x, y, w, h, r, b float32) {
	return body.x, body.y, body.w, body.h, body.x + body.w, body.y + body.h
}

func crossProduct(x1, y1, x2, y2 float32) float32 {
	return x1*y2 - y1*x2
}

func nearest(x, a, b float32) float32 {
	if abs(a-x) < abs(b-x) {
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
