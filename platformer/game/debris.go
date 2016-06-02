package game

import (
	"github.com/tanema/amore-examples/platformer/lib/bump"

	"github.com/tanema/amore/gfx"
)

type Debris struct {
	*Entity
	*bump.Body
	world    *bump.World
	color    *gfx.Color
	lifetime float32
	lived    float32
	vx, vy   float32
}

func newDebris(world *bump.World, x, y, r, g, b float32) *Debris {
	debris := &Debris{
		Entity:   newEntity(x, y, randRange(5, 10), randRange(5, 10)),
		world:    world,
		color:    gfx.NewColor(r, g, b, 255),
		lifetime: randRange(1, 4),
		lived:    0,
		vx:       randLimits(100),
		vy:       randLimits(100),
	}
	debris.Body = world.Add(
		debris, "debris",
		x, y, randRange(5, 10), randRange(5, 10),
		map[string]string{"block": "bounce", "guardian": "bounce"})
	return debris
}

func (debris *Debris) moveColliding(dt float32) {
	future_l := debris.l + debris.vx*dt
	future_t := debris.t + debris.vy*dt
	next_l, next_t, cols := debris.world.Move(future_l, future_t)

	for _, col := range cols {
		debris.changeVelocityByCollisionNormal(col.Normal.X, col.Normal.Y, 0.1)
	}

	debris.l, debris.t = next_l, next_t
}

func (debris *Debris) Send(event string, args ...interface{}) {
}

func (debris *Debris) Update(dt float32) {
	debris.lived = debris.lived + dt

	if debris.lived >= debris.lifeTime {
		debris.world.Remove(debris)
	} else {
		debris.changeVelocityByGravity(dt)
		debris.moveColliding(dt)
	}
}

func (debris *Debris) Draw() {
	r, g, b := debris.color.RGB()
	gfx.SetColor(r, g, b, 1000)
	gfx.Rect(gfx.FILL, debris.l, debris.t, debris.w, debris.h)
	gfx.SetColorC(debris.color)
	gfx.Rect(gfx.LINE, debris.l, debris.t, debris.w, debris.h)
}
