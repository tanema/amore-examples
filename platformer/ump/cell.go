package ump

type Cell struct {
	bodies map[uint32]*Body
}

func (cell *Cell) enter(body *Body) {
	cell.bodies[body.ID] = body
	body.cells = append(body.cells, cell)
}

func (cell *Cell) leave(body *Body) {
	delete(cell.bodies, body.ID)
}
