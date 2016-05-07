package game

import (
	"math/rand"

	"github.com/tanema/amore/keyboard"
)

type Player struct {
	*Car
	position float32
}

func newPlayer(z float32) *Player {
	return &Player{
		Car: newCar(segments[0], spriteSheet["player_straight"], z, 0, 0),
	}
}

func (player *Player) update(dt float32) {
	player.segment = findSegment(player.position + player.point.z)
	speedPercent := player.speed / maxSpeed
	dx := dt * 2 * speedPercent // at top speed, should be able to cross from left to right (-1 to 1) in 1 second

	player.position = increase(player.position, dt*player.speed, trackLength)

	updown := player.segment.p2.world.y - player.segment.p1.world.y
	if keyboard.IsDown(keyboard.KeyLeft) {
		player.point.x -= dx
		if updown > 0 {
			player.source = spriteSheet["player_uphill_left"]
		} else {
			player.source = spriteSheet["player_left"]
		}
	} else if keyboard.IsDown(keyboard.KeyRight) {
		player.point.x += dx
		if updown > 0 {
			player.source = spriteSheet["player_uphill_right"]
		} else {
			player.source = spriteSheet["player_right"]
		}
	} else {
		if updown > 0 {
			player.source = spriteSheet["player_uphill_straight"]
		} else {
			player.source = spriteSheet["player_straight"]
		}
	}

	player.point.x -= (dx * speedPercent * player.segment.curve * centrifugal)

	if keyboard.IsDown(keyboard.KeyUp) {
		player.speed = accelerate(player.speed, accel, dt)
	} else if keyboard.IsDown(keyboard.KeyDown) {
		player.speed = accelerate(player.speed, breaking, dt)
	} else {
		player.speed = accelerate(player.speed, decel, dt)
	}

	if (player.point.x < -1) || (player.point.x > 1) {
		if player.speed > offRoadLimit {
			player.speed = accelerate(player.speed, offRoadDecel, dt)
		}

		for n := 0; n < len(player.segment.sprites); n++ {
			sprite := player.segment.sprites[n]
			spriteW := sprite.source.GetWidth() * sprite_scale
			ov := float32(-1)
			if sprite.offset > 0 {
				ov = 1
			}
			if overlap(player.point.x, player.point.w, sprite.offset+spriteW/2*ov, spriteW, 1) {
				player.speed = maxSpeed / 5
				player.position = increase(player.segment.p1.world.z, -player.point.z, trackLength) // stop in front of sprite (at front of segment)
				break
			}
		}
	}

	for n := 0; n < len(player.segment.cars); n++ {
		car := player.segment.cars[n]
		if player.speed > car.speed {
			if overlap(player.point.x, player.point.w, car.offset, car.point.w, 0.8) {
				player.speed = car.speed * (car.speed / player.speed)
				player.position = increase(car.point.z, -player.point.z, trackLength)
				break
			}
		}
	}

	player.point.x = clamp(player.point.x, -3, 3)   // dont ever let it go too far out of bounds
	player.speed = clamp(player.speed, 0, maxSpeed) // or exceed maxSpeed
}

func (player *Player) draw() {
	playerPercent := percentRemaining(player.position+player.point.z, segmentLength)
	destY := (height / 2) - (cameraDepth / player.point.z * interpolate(player.segment.p1.camera.y, player.segment.p2.camera.y, playerPercent) * height / 2)
	side := float32(-1)
	if rand.Intn(1) == 0 {
		side = 1
	}
	bounce := (1.5 * rand.Float32() * player.speed / maxSpeed * resolution) * side
	player.Sprite.draw(width, height, roadWidth, cameraDepth/player.point.z, width/2, destY+bounce, -0.5, -1)
}
