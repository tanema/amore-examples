package game

import (
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"
)

const (
	playerAcc           = 200
	playerMaxSpeed      = 400
	playerRotationSpeed = 6
	playerFireRate      = 0.40
	playerJetSize       = 25
	playerJetWidth      = 0.15
)

type Player struct {
	*Sprite
	lastFire       float32
	isAccelerating bool
}

func newPlayer() *Player {
	new_player := &Player{}
	new_player.Sprite = NewSprite(new_player, "ship", screenWidth/2, screenHeight/2, 1,
		[]float32{
			-5, 4,
			0, -12,
			5, 4,
			-5, 4,
		}, true)
	return new_player
}

func (player *Player) Update(dt float32) {
	player.isAccelerating = false

	if keyboard.IsDown(keyboard.KeyLeft) {
		player.vrot = -playerRotationSpeed
	} else if keyboard.IsDown(keyboard.KeyRight) {
		player.vrot = playerRotationSpeed
	} else {
		player.vrot = 0
	}

	if keyboard.IsDown(keyboard.KeyUp) {
		player.isAccelerating = true
		player.ay = -(playerAcc * cos(player.rot))
		player.ax = playerAcc * sin(player.rot)
	} else {
		player.ax = 0
		player.ay = 0
	}

	player.lastFire += dt
	if keyboard.IsDown(keyboard.KeySpace) && player.lastFire > playerFireRate {
		addObject(newBullet(player.x, player.y, player.rot))
		lazer.Play()
		player.lastFire = 0
	}

	if collisions := player.UpdateMovement(dt); len(collisions) > 0 {
		for _, c := range collisions {
			if c.Name == "asteroid" {
				player.Destroy(false)
			}
		}
	}

	// limit the ship's speed
	if sqrt(player.vx*player.vx+player.vy*player.vy) > playerMaxSpeed {
		player.vx *= 0.95
		player.vy *= 0.95
	}
}

func (player *Player) Draw() {
	player.Sprite.Draw()

	if player.isAccelerating {
		points := player.Sprite.body.GetPoints()
		gfx.PolyLine([]float32{
			points[0] + ((points[4] - points[0]) * playerJetWidth),
			points[1] + ((points[5] - points[1]) * playerJetWidth),

			points[2] + (-sin(player.rot) * playerJetSize),
			points[3] + (cos(player.rot) * playerJetSize),

			points[4] + ((points[0] - points[4]) * playerJetWidth),
			points[5] + ((points[1] - points[5]) * playerJetWidth),
		})
	}
}

func (player *Player) Destroy(force bool) {
	removeObject(player)
	player.Sprite.Destroy()
	if !force {
		bomb.Play()
		gameOver = true
		newExplosion(player.Sprite)
	}
}
