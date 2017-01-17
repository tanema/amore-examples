package game

type Block struct {
	*Entity
	indestructible bool
}

func newBlock(gameMap *Map, l, t, w, h float32, indestructible bool) *Block {
	block := &Block{indestructible: indestructible}
	block.Entity = newEntity(gameMap, block, "block", l, t, w, h)
	block.body.SetStatic(true)
	return block
}

func (block *Block) getColor() (r, g, b float32) {
	if block.indestructible {
		return 150, 150, 220
	}
	return 220, 150, 150
}

func (block *Block) Update(dt float32) {
}

func (block *Block) Draw(debug bool) {
	r, g, b := block.getColor()
	l, t, w, h := block.Extents()
	drawFilledRectangle(l, t, w, h, r, g, b)
}

func (block *Block) damage(intensity float32) {
	if !block.indestructible {
		block.Entity.damage(intensity)
		area := block.w * block.h
		debrisNumber := floor(max(30, area/100))

		for i := float32(1); i <= debrisNumber; i++ {
			newDebris(block.gameMap,
				randRange(block.l, block.l+block.w),
				randRange(block.t, block.t+block.h),
				220, 150, 150,
			)
		}
	}
}
