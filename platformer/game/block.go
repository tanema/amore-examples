package game

import (
	"github.com/tanema/amore-examples/platformer/lib/bump"
)

type Block struct {
	*Entity
}

func newBlock(world *bump.World, l, t, w, h float32) *Block {
	newBlock := &Block{
		Entity: newEntity(world, l, t, w, h),
	}
	newBlock.Body = world.Add(newBlock, "block", l, t, w, h, map[string]string{})
	return newBlock
}

func (block *Block) Send(event string, args ...interface{}) {
}

func (block *Block) Update(dt float32) {
}
