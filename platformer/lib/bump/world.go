package bump

import (
	"math"
	"sort"

	"github.com/tanema/amore/gfx"
)

const defaultFilter = "slide"

type (
	World struct {
		cellSize  float32
		rows      map[int]map[int]*Cell
		bodies    map[Entity]*Body
		responses map[string]Resp
	}
	Resp func(col *Collision, body *Body, goalX, goalY float32) (gx, gy float32, cols []*Collision)
)

func NewWorld(cellSize int) *World {
	world := &World{
		cellSize:  float32(cellSize),
		rows:      make(map[int]map[int]*Cell),
		responses: map[string]Resp{},
		bodies:    map[Entity]*Body{},
	}

	world.AddResponse("touch", world.touch)
	world.AddResponse("cross", world.cross)
	world.AddResponse("slide", world.slide)
	world.AddResponse("bounce", world.bounce)

	return world
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
		if _, ok := visited[body]; !ok {
			visited[body] = true
			ok, ti1, ti2, _, _, _, _ := body.getSegmentIntersectionIndices(x1, y1, x2, y2, 0, 1)
			if ok && ((0 < ti1 && ti1 < 1) || (0 < ti2 && ti2 < 1)) {
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

func (world *World) Add(entity Entity, tag string, left, top, width, height float32, respMap map[string]string) *Body {
	body := &Body{
		Entity:  entity,
		world:   world,
		tag:     tag,
		width:   width,
		height:  height,
		cells:   []*Cell{},
		respMap: respMap,
	}
	world.bodies[entity] = body
	body.move(left, top)
	return body
}

func (world *World) Remove(entity Entity) {
	body := world.bodies[entity]
	body.remove()
	delete(world.bodies, entity)
}

func (world *World) Update(entity Entity, x, y, w, h float32) {
	body := world.bodies[entity]
	body.update(x, y, w, h)
}

func (world *World) Move(entity Entity, goalX, goalY float32) (gx, gy float32, cols []*Collision) {
	body := world.bodies[entity]
	return body.move(goalX, goalY)
}

func (world *World) Check(entity Entity, goalX, goalY float32) (gx, gy float32, cols []*Collision) {
	body := world.bodies[entity]
	return body.check(goalX, goalY)
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

func (world *World) getEntitiesInCells(cells []*Cell) []Entity {
	dict := make(map[*Body]bool)
	entities := []Entity{}
	for _, cell := range cells {
		for _, body := range cell.bodies {
			if _, ok := dict[body]; !ok {
				entities = append(entities, body.Entity)
				dict[body] = true
			}
		}
	}
	return entities
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
	tr := float32(math.Max(float64(goalX+body.width), float64(body.x+body.width)))
	tb := float32(math.Max(float64(goalY+body.height), float64(body.y+body.height)))

	tw, th := tr-tl, tb-tt

	visited := map[*Body]bool{}

	bodies := world.getBodiesInCells(world.cellsInRect(tl, tt, tw, th))
	for _, other := range bodies {
		if _, ok := visited[body]; !ok {
			visited[body] = true

			if col := body.collide(other, goalX, goalY); col != nil {
				col.Body = body
				col.Other = other

				var ok bool
				col.RespType, ok = body.respMap[other.tag]
				if !ok {
					col.RespType = defaultFilter
				}

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

func (world *World) touch(col *Collision, body *Body, goalX, goalY float32) (float32, float32, []*Collision) {
	return col.Touch.X, col.Touch.Y, []*Collision{}
}

func (world *World) cross(col *Collision, body *Body, goalX, goalY float32) (float32, float32, []*Collision) {
	return goalX, goalY, world.Project(body, goalX, goalY)
}

func (world *World) slide(col *Collision, body *Body, goalX, goalY float32) (float32, float32, []*Collision) {
	sx, sy := col.Touch.X, col.Touch.Y
	if col.Move.X != 0 || col.Move.Y != 0 {
		if col.Normal.X == 0 {
			sx = goalX
		} else {
			sy = goalY
		}
	}

	col.Data = Point{X: sx, Y: sy}

	return sx, sy, world.Project(&Body{
		x:      col.Touch.X,
		y:      col.Touch.Y,
		width:  body.width,
		height: body.height,
	}, sx, sy)
}

func (world *World) bounce(col *Collision, body *Body, goalX, goalY float32) (float32, float32, []*Collision) {
	tx, ty := col.Touch.X, col.Touch.Y
	bx, by := tx, ty
	if col.Move.X != 0 || col.Move.Y != 0 {
		bnx, bny := goalX-tx, goalY-ty
		if col.Normal.X == 0 {
			bny = -bny
		} else {
			bnx = -bnx
		}
		bx, by = tx+bnx, ty+bny
	}
	col.Data = Point{X: bx, Y: by}
	body.x, body.y = col.Touch.X, col.Touch.Y
	return bx, by, world.Project(body, bx, by)
}

func (world *World) DrawDebug(l, t, w, h float32) {
	cl, ct, cw, ch := world.gridToCellRect(l, t, w, h)
	for cy := ct - 1; cy <= ct+ch-1; cy++ {
		for cx := cl - 1; cx <= cl+cw-1; cx++ {
			gfx.SetLineWidth(2)
			gfx.SetColor(255, 255, 255, 200)
			gfx.Rect(gfx.LINE, float32(cx)*world.cellSize, float32(cy)*world.cellSize, world.cellSize, world.cellSize)
		}
	}

	for _, cell := range world.cellsInRect(l, t, w, h) {
		cell.DrawDebug(world.cellSize)
	}
}
