package ump

type Cell struct {
	bodies map[uint32]*Body
}

func (cell *Cell) enter(body *Body) {
	cell.bodies[body.ID] = body
}

func (cell *Cell) leave(body *Body) {
	delete(cell.bodies, body.ID)
}
