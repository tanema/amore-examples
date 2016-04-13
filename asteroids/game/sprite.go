package game

import (
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"

	"github.com/tanema/amore-examples/asteroids/game/phys"
)

type Sprite struct {
	name         string
	body         *phys.Body
	scale        float32
	x, y, rot    float32
	vx, vy, vrot float32
	ax, ay       float32
	wraps        bool
}

func NewSprite(collidable phys.Collidable, name string, x, y, scale float32, points []float32, wraps bool) *Sprite {
	new_sprite := &Sprite{
		name:  name,
		body:  world.AddBody(collidable, name, x, y, scale, points),
		x:     x,
		y:     y,
		scale: scale,
		wraps: wraps,
	}
	return new_sprite
}

func (sprite *Sprite) UpdateMovement(delta float32) []*phys.Body {
	sprite.vx += sprite.ax * delta
	sprite.vy += sprite.ay * delta
	dx, dy, dr := sprite.vx*delta, sprite.vy*delta, sprite.vrot*delta
	sprite.x, sprite.y, sprite.rot = sprite.x+dx, sprite.y+dy, sprite.rot+dr
	collisions := sprite.body.Move(sprite.x, sprite.y, sprite.rot, sprite.scale)
	if sprite.wraps {
		if sprite.x > screenWidth {
			sprite.x = 0
		} else if sprite.x < 0 {
			sprite.x = screenWidth
		}
		if sprite.y > screenHeight {
			sprite.y = 0
		} else if sprite.y < 0 {
			sprite.y = screenHeight
		}
	}
	return collisions
}

func (sprite *Sprite) Draw() {
	if keyboard.IsDown(keyboard.KeyG) {
		sprite.body.Draw()
	}
	gfx.SetLineWidth(2)
	gfx.PolyLine(sprite.body.GetPoints())
}

func (sprite *Sprite) GetPoints() []float32 {
	return sprite.body.GetPoints()
}

func (sprite *Sprite) Destroy() {
	sprite.body.Remove()
}
