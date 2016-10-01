package game

import (
	"github.com/tanema/amore-examples/platformer/lib/bump"
	"github.com/tanema/amore/gfx"
)

type Block struct {
	x, y, width, height float32
	color               *gfx.Color
	body                *bump.Body
}

func NewBlock(x, y, width, height float32, color *gfx.Color) *Block {
	newBlock := &Block{
		x:      x,
		y:      y,
		width:  width,
		height: height,
		color:  color,
	}
	newBlock.body = world.Add(newBlock, "block", x, y, width, height, map[string]string{})
	return newBlock
}

func (block *Block) Send(event string, args ...interface{}) {
}

func (block *Block) Update(dt float32) {
	block.x, block.y = block.body.Position()
}

func (block *Block) Draw() {
	gfx.SetColorC(block.color)
	gfx.Rect(gfx.FILL, block.x, block.y, block.width, block.height)
}
