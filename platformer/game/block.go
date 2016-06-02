package game

import (
	"math"

	"github.com/tanema/amore-examples/platformer/lib/bump"

	"github.com/tanema/amore/gfx"
)

type Block struct {
	*Entity
	*bump.Body
	world          *bump.World
	indestructible bool
	color          *gfx.Color
}

func newBlock(world *bump.World, l, t, w, h float32, indestructible bool) *Block {
	color := gfx.NewColor(220, 150, 150, 255)
	if indestructible {
		color = gfx.NewColor(150, 150, 220, 255)
	}
	newBlock := &Block{
		world:          world,
		Entity:         newEntity(l, t, w, h),
		indestructible: indestructible,
		color:          color,
	}
	newBlock.Body = world.Add(newBlock, "block", l, t, w, h, map[string]string{})
	return newBlock
}

func (block *Block) Draw() {
	gfx.SetColor(block.color[0], block.color[1], block.color[2], 1000)
	gfx.Rect(gfx.FILL, block.l, block.t, block.w, block.h)
	gfx.SetColorC(block.color)
	gfx.Rect(gfx.LINE, block.l, block.t, block.w, block.h)
}

func (block *Block) Send(event string, args ...interface{}) {
}

func (block *Block) Update(dt float32) {
}

func (block *Block) destroy() {
	block.world.Remove(block)
	area := float64(block.w * block.h)
	debrisNumber := int(math.Floor(math.Max(30, area/100)))
	for i := 1; i <= debrisNumber; i++ {
		newDebris(block.world,
			randRange(block.l, block.l+block.w),
			randRange(block.t, block.t+block.h),
			220, 150, 150)
	}
}
