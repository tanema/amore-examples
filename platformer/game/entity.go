package game

const gravityAccel float32 = 500 // pixels per second^2

type Entity struct {
	l, t, w, h float32
	vx, vy     float32
}

func newEntity(l, t, w, h float32) *Entity {
	return &Entity{l: l, t: t, w: w, h: h}
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
