package game

import (
	"github.com/tanema/amore/gfx"
)

type Grenade struct {
	*Entity
	parent        *Guardian
	lived         float32
	ignoresParent bool
}

const (
	grenadeLifeTime   = float32(4)
	grenadeBounciness = float32(0.4)
)

func newGrenade(guardian *Guardian, x, y, vx, vy float32) *Grenade {
	grenade := &Grenade{
		parent:        guardian,
		ignoresParent: true,
	}
	grenade.Entity = newEntity(guardian.gameMap, grenade, "grenade", x, y, 11, 11)
	grenade.vx, grenade.vy = vx, vy
	grenade.body.SetResponses(map[string]string{
		"guardian": "cross", //This gets toggled once the grenade is outside the guardian
		"player":   "bounce",
		"block":    "bounce",
	})
	return grenade
}

func (grenade *Grenade) moveColliding(dt float32) {
	future_l := grenade.l + grenade.vx*dt
	future_t := grenade.t + grenade.vy*dt
	next_l, next_t, cols := grenade.body.Move(future_l, future_t)

	for _, col := range cols {
		if col.Body.Tag() == "player" {
			grenade.destroy()
			return
		}
		grenade.changeVelocityByCollisionNormal(col.Normal.X, col.Normal.Y, grenadeBounciness)
	}
	grenade.l, grenade.t = next_l, next_t
}

func (grenade *Grenade) detectExitedParent() {
	if grenade.ignoresParent {
		x1, y1, w1, h1 := grenade.Extents()
		x2, y2, w2, h2 := grenade.parent.Extents()
		grenade.ignoresParent = x1 < x2+w2 && x2 < x1+w1 && y1 < y2+h2 && y2 < y1+h1
		if !grenade.ignoresParent {
			grenade.body.SetResponse("guardian", "bounce")
		}
	}
}

func (grenade *Grenade) updateOrder() int {
	return 2
}

func (grenade *Grenade) update(dt float32) {
	grenade.lived += dt
	if grenade.lived >= grenadeLifeTime {
		grenade.destroy()
	} else {
		grenade.changeVelocityByGravity(dt)
		grenade.moveColliding(dt)
		grenade.detectExitedParent()
	}
}

func (grenade *Grenade) draw(debug bool) {
	r, g, b := float32(255), float32(0), float32(0)
	gfx.SetColor(r, g, b, 255)

	cx, cy := grenade.GetCenter()
	gfx.Circle(gfx.LINE, cx, cy, 8)

	percent := grenade.lived / grenadeLifeTime
	g = floor(255 * percent)
	b = g

	gfx.SetColor(r, g, b, 255)
	gfx.Circle(gfx.FILL, cx, cy, 8)

	if debug {
		gfx.SetColor(255, 255, 255, 200)
		l, t, w, h := grenade.Extents()
		gfx.Rect(gfx.LINE, l, t, w, h)
	}
}

func (grenade *Grenade) destroy() {
	grenade.Entity.destroy()
	newExplosion(grenade)
}
