package game

import (
	"math"
)

const (
	explosionSpeed = 100
	explosionTime  = 3
)

type Explosion struct {
	sprites []*Sprite
	life    float32
}

func newExplosion(sprite *Sprite) {
	points := sprite.GetPoints()
	explosion := &Explosion{
		sprites: []*Sprite{},
	}

	for i := 0; i < len(points)-2; i += 2 {
		explosion.addSegment(points[i], points[i+1], points[i+2], points[i+3])
	}

	addObject(explosion)
}

func (explosion *Explosion) addSegment(x0, y0, x1, y1 float32) {
	sprite := NewSprite(nil, "explosion", 0, 0, 1, []float32{x0, y0, x1, y1}, false)
	rot := math.Atan2(float64(x1-x0), float64(y1-y0))
	vectorx := float32(math.Sin(float64(rot)))
	vectory := -float32(math.Cos(float64(rot)))
	sprite.vx = explosionSpeed * vectorx
	sprite.vy = explosionSpeed * vectory
	explosion.sprites = append(explosion.sprites, sprite)
}

func (explosion *Explosion) Update(dt float32) {
	for _, sprite := range explosion.sprites {
		sprite.UpdateMovement(dt)
	}
	explosion.life += dt
	if explosion.life >= explosionTime {
		explosion.Destroy(false)
	}
}

func (explosion *Explosion) Draw() {
	for _, sprite := range explosion.sprites {
		sprite.Draw()
	}
}

func (explosion *Explosion) Destroy(force bool) {
	removeObject(explosion)
	for _, sprite := range explosion.sprites {
		sprite.Destroy()
	}
}
