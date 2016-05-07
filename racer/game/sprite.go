package game

import (
	"github.com/tanema/amore/gfx"
)

type Sprite struct {
	source *gfx.Quad
	offset float32
}

func newSprite(source *gfx.Quad, offset float32) *Sprite {
	return &Sprite{
		source: source,
		offset: offset,
	}
}

func (sprite *Sprite) draw(width, height, roadWidth, scale, destX, destY, offsetX, offsetY float32) {
	w, h := float32(sprite.source.GetWidth()), float32(sprite.source.GetHeight())
	destW := (w * scale * width / 2) * (sprite_scale * roadWidth)
	destH := (h * scale * width / 2) * (sprite_scale * roadWidth)
	gfx.SetColor(255, 255, 255, 255)
	gfx.Drawq(sprites, sprite.source, destX+(destW*offsetX), destY+(destH*offsetY), 0, destW/w, destH/h)
}
