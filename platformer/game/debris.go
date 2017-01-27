package game

type Debris struct {
	*Entity
	r, g, b  float32
	lifeTime float32
	lived    float32
}

func newDebris(gameMap *Map, x, y, r, g, b float32) *Debris {
	debris := &Debris{
		r:        r,
		g:        g,
		b:        b,
		lifeTime: 1 + 3*randMax(1),
	}
	debris.Entity = newEntity(gameMap,
		debris, "debris",
		x, y,
		randRange(5, 10),
		randRange(5, 10),
	)
	debris.vx = randRange(-50, 50)
	debris.vy = randRange(-50, 50)
	debris.body.SetResponses(map[string]string{
		"guardian": "bounce",
		"block":    "bounce",
	})
	return debris
}

func (debris *Debris) moveColliding(dt float32) {
	future_l := debris.l + debris.vx*dt
	future_t := debris.t + debris.vy*dt
	next_l, next_t, cols := debris.Entity.body.Move(future_l, future_t)
	for _, col := range cols {
		debris.changeVelocityByCollisionNormal(col.Normal.X, col.Normal.Y, 0.1)
	}
	debris.l, debris.t = next_l, next_t
}

func (debris *Debris) update(dt float32) {
	debris.lived = debris.lived + dt

	if debris.lived >= debris.lifeTime {
		debris.destroy()
	} else {
		debris.changeVelocityByGravity(dt)
		debris.moveColliding(dt)
	}
}

func (debris *Debris) draw(debug bool) {
	l, t, w, h := debris.Extents()
	drawFilledRectangle(l, t, w, h, debris.r, debris.g, debris.b)
}
