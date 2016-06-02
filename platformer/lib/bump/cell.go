package bump

type Cell struct {
	world  *World
	cx     int
	cy     int
	bodies []*Body
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
