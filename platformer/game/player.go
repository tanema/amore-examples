package game

import (
	"github.com/tanema/amore/keyboard"
)

type Player struct {
	*Entity
	health             float32
	deadCounter        float32
	isJumpingOrFlying  bool
	isDead             bool
	onGround           bool
	achievedFullHealth bool
}

const (
	deadDuration float32 = 3   // seconds until res-pawn
	runAccel     float32 = 500 // the player acceleration while going left/right
	brakeAccel   float32 = 2000
	jumpVelocity float32 = 400 // the initial upwards velocity when jumping
	beltWidth    float32 = 2
	beltHeight   float32 = 8
)

func newPlayer(gameMap *Map, l, t float32) *Player {
	player := &Player{
		health: 1,
	}
	player.Entity = newEntity(gameMap, player, "player", l, t, 32, 64)
	player.body.SetResponses(map[string]string{
		"guardian": "slide",
		"block":    "slide",
	})
	return player
}

func (player *Player) changeVelocityByKeys(dt float32) {
	player.isJumpingOrFlying = false

	if player.isDead {
		return
	}

	if keyboard.IsDown(keyboard.KeyLeft) {
		if player.vx > 0 {
			player.vx -= dt * brakeAccel
		} else {
			player.vx -= dt * runAccel
		}
	} else if keyboard.IsDown(keyboard.KeyRight) {
		if player.vx < 0 {
			player.vx += dt * brakeAccel
		} else {
			player.vx += dt * runAccel
		}
	} else {
		brake := dt * -brakeAccel
		if player.vx < 0 {
			brake = dt * brakeAccel
		}
		if abs(brake) > abs(player.vx) {
			player.vx = 0
		} else {
			player.vx += brake
		}
	}

	if keyboard.IsDown(keyboard.KeyUp) && (player.canFly() || player.onGround) { // jump/fly
		player.vy = -jumpVelocity
		player.isJumpingOrFlying = true
	}
}

func (player *Player) moveColliding(dt float32) {
	player.onGround = false
	l, t, cols := player.Entity.body.Move(player.l+player.vx*dt, player.t+player.vy*dt)
	for _, col := range cols {
		if col.Body.Tag() != "puff" {
			player.changeVelocityByCollisionNormal(col.Normal.X, col.Normal.Y, 0)
			player.onGround = col.Normal.Y < 1
		}
	}
	player.l, player.t = l, t
}

func (player *Player) updateHealth(dt float32) {
	player.achievedFullHealth = false
	if player.health < 1 {
		player.health = min(1, player.health+dt/6)
		player.achievedFullHealth = player.health == 1
	}
}

func (player *Player) playEffects() {
	if player.isJumpingOrFlying {
		if !player.onGround {
			l, t, w, h := player.Extents()
			newPuff(player.gameMap, l, t+h/2, 20*(1-randMax(1)), 50, 2, 3)
			newPuff(player.gameMap, l+w, t+h/2, 20*(1-randMax(1)), 50, 2, 3)
		}
	}
}

func (player *Player) updateOrder() int {
	return 1
}

func (player *Player) update(dt float32) {
	player.updateHealth(dt)
	player.changeVelocityByKeys(dt)
	player.changeVelocityByGravity(dt)
	player.playEffects()
	player.moveColliding(dt)
}

func (player *Player) getColor() (r, g, b float32) {
	g = floor(255 * player.health)
	return 255 - g, g, 0
}

func (player *Player) canFly() bool {
	return player.health == 1
}

func (player *Player) draw(debug bool) {
	r, g, b := player.getColor()
	l, t, w, h := player.Extents()
	drawFilledRectangle(l, t, w, h, r, g, b)

	if player.canFly() {
		drawFilledRectangle(l-beltWidth, t+h/2, w+2*beltWidth, beltHeight, 255, 255, 255)
	}

	if debug && player.onGround {
		drawFilledRectangle(l, t+h-4, w, 4, 255, 255, 255)
	}
}

func (player *Player) damage(intensity float32) {
	if player.isDead {
		return
	}

	if player.health == 1 {
		for i := 1; i <= 3; i++ {
			newDebris(player.gameMap,
				randRange(player.l, player.l+player.w),
				player.t+player.h/2,
				255, 0, 0,
			)
		}
	}

	player.health = player.health - intensity
	if player.health <= 0 {
		player.destroy()
		player.isDead = true
	}
}

func (player *Player) destroy() {
	player.body.Remove()
	for i := 1; i <= 20; i++ {
		newDebris(player.gameMap,
			randRange(player.l, player.l+player.w),
			randRange(player.t, player.t+player.h),
			255, 0, 0)
	}
}
