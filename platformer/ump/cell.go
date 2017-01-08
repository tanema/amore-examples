package ump

type Cell struct {
	bodies    map[uint32]*Body
	itemCount int
}

func (cell *Cell) enter(body *Body) {
	cell.bodies[body.ID] = body
	body.cells = append(body.cells, cell)
	cell.itemCount = len(cell.bodies)
}

func (cell *Cell) leave(body *Body) {
	delete(cell.bodies, body.ID)
	cell.itemCount = len(cell.bodies)
}
