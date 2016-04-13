package phys

type Cell struct {
	x, y, width, height float32
	bodies              []*Body
}

func (cell *Cell) enter(body *Body) {
	cell.leave(body) // make sure this body is uniqe
	cell.bodies = append(cell.bodies, body)
}

func (cell *Cell) leave(body *Body) {
	for i, b := range cell.bodies {
		if body == b {
			cell.bodies = append(cell.bodies[:i], cell.bodies[i+1:]...)
			return
		}
	}
}
