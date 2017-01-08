package game

import (
	"github.com/tanema/amore/timer"

	"github.com/tanema/amore-examples/platformer/ump"
)

const gravityAccel float32 = 500 // pixels per second^2

type Entity struct {
	body_tag   string
	l, t, w, h float32
	vx, vy     float32
	gameMap    *Map
	body       *ump.Body
	created_at float32
}

func newEntity(gameMap *Map, object gameObject, tag string, l, t, w, h float32) *Entity {
	entity := &Entity{
		body_tag: tag,
		gameMap:  gameMap,
		l:        l, t: t, w: w, h: h,
		created_at: timer.GetTime(),
	}
	entity.body = gameMap.world.Add(tag, l, t, w, h)
	gameMap.objects[entity.body.ID] = object
	return entity
}

func (entity *Entity) changeVelocityByGravity(dt float32) {
	entity.vy += gravityAccel * dt
}

func (entity *Entity) changeVelocityByCollisionNormal(nx, ny, bounciness float32) {
	if (nx < 0 && entity.vx > 0) || (nx > 0 && entity.vx < 0) {
		entity.vx = -entity.vx * bounciness
	}

	if (ny < 0 && entity.vy > 0) || (ny > 0 && entity.vy < 0) {
		entity.vy = -entity.vy * bounciness
	}
}

func (entity *Entity) GetCenter() (x, y float32) {
	return entity.l + (entity.w / 2), entity.t + (entity.h / 2)
}

func (entity *Entity) Extents() (l, t, w, h float32) {
	return entity.l, entity.t, entity.w, entity.h
}

func (entity *Entity) destroy() {
	entity.body.Remove()
	//delete(entity.gameMap.objects, entity.body.ID)
}

func (entity *Entity) tag() string {
	return entity.body_tag
}

func (entity *Entity) push(strength float32) {
	cx, cy := entity.l+entity.w/2, entity.t+entity.h/2
	ox, oy := entity.GetCenter()
	dx, dy := ox-cx, oy-cy
	dx, dy = clamp(dx, -strength, strength), clamp(dy, -strength, strength)

	if dx > 0 {
		dx = strength - dx
	} else if dx < 0 {
		dx = dx - strength
	}

	if dy > 0 {
		dy = strength - dy
	} else if dy < 0 {
		dy = dy - strength
	}

	entity.vx += dx
	entity.vy += dy
}

func (entity *Entity) damage(intensity float32) {
	entity.destroy()
}
