package game

func generateTerrain(size, iterations int, smooth bool) [][]*Cell {
	terrain := make([][]*Cell, size)
	for x := 0; x < size; x++ {
		terrain[x] = make([]*Cell, size)
		for y := 0; y < size; y++ {
			terrain[x][y] = newCell(x, y)
		}
	}

	for ; iterations >= 0; iterations-- {
		px, py, r := randRange(0, size), randRange(0, size), randRange(10, 40)
		for x := -r; x <= r; x++ {
			for y := -r; y <= r; y++ {
				// Increase altitude of cell cell with "bell" function factor
				cell := terrain[(px+x+size)%size][(py-y+size)%size]
				cell.setZ(cell.getZ() + 5*exp(-(float32(x)*float32(x)+float32(y)*float32(y))/(float32(r)*2)))
			}
		}
	}

	if !smooth {
		for x := 0; x < size; x++ {
			for y := 0; y < size; y++ {
				terrain[x][y].setZ(floor(terrain[x][y].getZ()))
			}
		}
	}

	return terrain
}
