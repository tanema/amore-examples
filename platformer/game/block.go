package game

import (
	"github.com/tanema/amore-examples/platformer/ump"
	"github.com/tanema/amore/gfx"
)

type Block struct {
	x, y, width, height float32
	color               *gfx.Color
	body                *ump.Body
}

func NewBlock(x, y, width, height float32, color *gfx.Color) *Block {
	newBlock := &Block{
		x:      x,
		y:      y,
		width:  width,
		height: height,
		color:  color,
	}
	newBlock.body = world.Add("block", x, y, width, height)
	newBlock.x, newBlock.y = newBlock.body.Position()
	return newBlock
}

func (block *Block) Update(dt float32) {
	block.x, block.y = block.body.Position()
}

func (block *Block) Draw() {
	gfx.SetColorC(block.color)
	gfx.Rect(gfx.FILL, block.x, block.y, block.width, block.height)
}
