package game

import (
	"github.com/tanema/amore/gfx"
)

type Guardian struct {
	*Entity
	gameMap                    *Map
	activeRadius               float32
	fireCoolDown               float32 // how much time the Guardian takes to "regenerate a grenade"
	aimDuration                float32 // time it takes to "aim"
	targetCoolDown             float32 // minimum time between "target acquired" chirps
	fireTimer                  float32
	aimTimer                   float32
	timeSinceLastTargetAquired float32
	isNearTarget               bool
	isLoading                  bool
	laserX, laserY             float32
}

func newGuardian(gameMap *Map, l, t float32) *Guardian {
	guardian := &Guardian{
		gameMap:                    gameMap,
		activeRadius:               500,
		fireCoolDown:               0.75,
		aimDuration:                1.25,
		targetCoolDown:             2,
		timeSinceLastTargetAquired: 2,
	}
	guardian.Entity = newEntity(gameMap, guardian, "guardian", l, t, 42, 110)

	l, t, w, h := guardian.Extents()
	others := gameMap.world.QueryRect(l, t, w, h, "block")
	for _, other := range others {
		if other.ID != guardian.body.ID {
			other.Remove()
		}
	}

	return guardian
}

func (guardian *Guardian) update(dt float32) {
	guardian.isNearTarget = false
	guardian.isLoading = false
	guardian.laserX, guardian.laserY = 0, 0
	guardian.timeSinceLastTargetAquired += dt

	if guardian.fireTimer < guardian.fireCoolDown {
		guardian.fireTimer = guardian.fireTimer + dt
		guardian.isLoading = true
	} else {
		cx, cy := guardian.GetCenter()
		tx, ty := guardian.gameMap.Player.GetCenter()
		dx, dy := cx-tx, cy-ty
		distance2 := dx*dx + dy*dy

		if distance2 <= guardian.activeRadius*guardian.activeRadius {
			guardian.isNearTarget = true
			bodies := guardian.gameMap.world.QuerySegment(cx, cy, tx, ty, "player", "block", "guardian")
			// ignore itemsInfo[0] because thats always guardian
			if len(bodies) > 1 {
				body := bodies[1]
				if body.ID == guardian.gameMap.Player.body.ID {
					guardian.laserX, guardian.laserY = guardian.gameMap.Player.GetCenter()
					if guardian.aimTimer == 0 && guardian.timeSinceLastTargetAquired >= guardian.targetCoolDown {
						guardian.timeSinceLastTargetAquired = 0
					}
					guardian.aimTimer = guardian.aimTimer + dt
					if guardian.aimTimer >= guardian.aimDuration {
						guardian.fire()
					}
				} else {
					guardian.aimTimer = 0
				}
			}
		} else {
			guardian.aimTimer = 0
		}
	}
}

func (guardian *Guardian) updateOrder() int {
	return 3
}

func (guardian *Guardian) draw(debug bool) {
	drawFilledRectangle(guardian.l, guardian.t, guardian.w, guardian.h, 255, 0, 255)

	cx, cy := guardian.GetCenter()
	gfx.SetColor(255, 0, 0, 255)
	radius := float32(8)
	if guardian.isLoading {
		percent := guardian.fireTimer / guardian.fireCoolDown
		alpha := floor(255 * percent)
		radius = radius * percent

		gfx.SetColor(0, 100, 200, alpha)
		gfx.Circle(gfx.FILL, cx, cy, radius)
		gfx.SetColor(0, 100, 200, 255)
		gfx.Circle(gfx.LINE, cx, cy, radius)
	} else {
		if guardian.aimTimer > 0 {
			gfx.SetColor(255, 0, 0, 255)
		} else {
			gfx.SetColor(0, 100, 200, 255)
		}
		gfx.Circle(gfx.LINE, cx, cy, radius)
		gfx.Circle(gfx.FILL, cx, cy, radius)

		if debug {
			gfx.SetColor(255, 255, 255, 100)
			gfx.Circle(gfx.LINE, cx, cy, guardian.activeRadius)
		}

		if guardian.isNearTarget {
			tx, ty := guardian.gameMap.Player.GetCenter()

			if debug {
				gfx.SetColor(255, 255, 255, 100)
				gfx.Line(cx, cy, tx, ty)
			}

			if guardian.aimTimer > 0 {
				gfx.SetColor(255, 100, 100, 200)
			} else {
				gfx.SetColor(0, 100, 200, 100)
			}
			if guardian.laserX != 0 && guardian.laserY != 0 {
				gfx.SetLineWidth(2)
				gfx.Line(cx, cy, guardian.laserX, guardian.laserY)
				gfx.SetLineWidth(1)
			}
		}
	}
}

func (guardian *Guardian) fire() {
	cx, cy := guardian.GetCenter()
	tx, ty := guardian.gameMap.Player.GetCenter()
	vx, vy := (tx-cx)*3, (ty-cy)*3
	newGrenade(guardian, cx, cy, vx, vy)
	guardian.fireTimer = 0
	guardian.aimTimer = 0
}

func (guardian *Guardian) damage(intensity float32) {
	guardian.destroy()
}

func (guardian *Guardian) destroy() {
	guardian.body.Remove()
	for i := 1; i <= 45; i++ {
		newDebris(guardian.gameMap,
			randRange(guardian.l, guardian.l+guardian.w),
			randRange(guardian.t, guardian.t+guardian.h),
			255, 0, 255,
		)
	}
}
