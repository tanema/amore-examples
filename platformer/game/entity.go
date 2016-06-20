package game

import (
	"github.com/tanema/amore-examples/platformer/lib/bump"

	"github.com/tanema/amore/gfx"
)

const gravityAccel float32 = 500 // pixels per second^2

type Entity struct {
	l, t, w, h float32
	vx, vy     float32
	*bump.Body
	world *bump.World
}

func newEntity(world *bump.World, l, t, w, h float32) *Entity {
	return &Entity{world: world, l: l, t: t, w: w, h: h}
}

func (entity *Entity) Draw() {
	gfx.SetColor(255, 255, 255, 255)
	gfx.Rect(gfx.LINE, entity.l, entity.t, entity.w, entity.h)
}

func (entity *Entity) changeVelocityByGravity(dt float32) {
	entity.vy = entity.vy + gravityAccel*dt
}

func (entity *Entity) changeVelocityByCollisionNormal(nx, ny, bounciness float32) {
	if (nx < 0 && entity.vx > 0) || (nx > 0 && entity.vx < 0) {
		entity.vx = -entity.vx * bounciness
	}
	if (ny < 0 && entity.vy > 0) || (ny > 0 && entity.vy < 0) {
		entity.vy = -entity.vy * bounciness
	}
}

func (entity *Entity) getCenter() (x, y float32) {
	return entity.l + entity.w/2, entity.t + entity.h/2
}
