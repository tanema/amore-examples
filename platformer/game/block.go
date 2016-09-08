package game

import (
	"github.com/tanema/amore/gfx"
)

type Block struct {
	x, y, width, height float32
	color               *gfx.Color
}

func NewBlock(x, y, width, height float32, color *gfx.Color) *Block {
	return &Block{x, y, width, height, color}
}

func (block *Block) Send(event string, args ...interface{}) {
}

func (block *Block) Update(dt float32) {
}

func (block *Block) Draw() {
	gfx.SetColorC(block.color)
	gfx.Rect(gfx.FILL, block.x, block.y, block.width, block.height)
}
