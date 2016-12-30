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

func (block *Block) Draw() {
	r, g, b := block.getColor()
	l, t, w, h := block.Extents()
	drawFilledRectangle(l, t, w, h, r, g, b)
}

func (block *Block) destroy() {
	block.Entity.destroy()
	//local area = self.w * self.h
	//local debrisNumber = math.floor(math.max(30, area / 100))

	//for i=1, debrisNumber do
	//Debris:new(self.world,
	//math.random(self.l, self.l + self.w),
	//math.random(self.t, self.t + self.h),
	//220, 150, 150
	//)
	//end
}
