package game

import (
	"math"
)

const (
	bulletSpeed = 500
)

type Bullet struct {
	*Sprite
}

func newBullet(x, y, rot float32) *Bullet {
	vectorx := float32(math.Sin(float64(rot)))
	vectory := -float32(math.Cos(float64(rot)))

	bullet := &Bullet{}
	bullet.Sprite = NewSprite(bullet, "bullet", x+(vectorx*10), y+(vectory*10), 1,
		[]float32{
			-1, 0,
			1, 0,
		}, false)
	bullet.rot = rot
	bullet.vx = (bulletSpeed * vectorx)
	bullet.vy = (bulletSpeed * vectory)

	return bullet
}

func (bullet *Bullet) Update(dt float32) {
	if collisions := bullet.UpdateMovement(dt); len(collisions) > 0 {
		bullet.Destroy(false)
		for _, c := range collisions {
			if c.Name == "asteroid" {
				score++
				c.Collidable.Destroy(false)
			}
		}
	}

	if bullet.x > screenWidth || bullet.x < 0 || bullet.y > screenHeight || bullet.y < 0 {
		bullet.Destroy(false)
	}
}

func (bullet *Bullet) Destroy(force bool) {
	removeObject(bullet)
	bullet.Sprite.Destroy()
}
