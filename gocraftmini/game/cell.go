package game

// Cell describes one position on the grid
type Cell struct {
	biom *Voxel
	dirt *Voxel
}

func newCell(x, y int) *Cell {
	return &Cell{
		biom: newVoxel(float32(x), float32(y), -9, 1, 0, 0, 0, 0, true),
		dirt: newVoxel(float32(x), float32(y), 0, 1, 0, 0.083, worldSaturation, 0.2, true),
	}
}

func (cell *Cell) getZ() float32 {
	return cell.biom.z
}

func (cell *Cell) setZ(newZ float32) {
	cell.biom.z, cell.dirt.height = newZ, newZ
	var h, s, l float32
	if cell.biom.z < 0 { // water, display depth of water, deeper is darker
		h, s, l = (180-cell.biom.z*20/3)/360, 0.99, 0.5
	} else if cell.biom.z == 0 { // sand
		h, s, l = 0.16, worldSaturation, 0.6
	} else if cell.biom.z > 15 { // snow
		h, s, l = 0, 0, 0.99
	} else { // grass
		h, s, l = 0.33, worldSaturation, 0.3
	}
	cell.biom.h, cell.biom.s, cell.biom.l = h, s, l
}

func (cell *Cell) update(world *World, x, y, distX, distY float32) {
	cell.dirt.update(world, distX, distY)
	cell.biom.update(world, distX, distY)
	cell.biom.height = 0
	if cell.biom.z < 0 { // biom is water
		// waves in the water
		cell.biom.height = (1.4 + sin(world.timeOfDay*25+cell.biom.y)) / 6 // waves
		// reflection off of the water
		sunx := float32(world.camera.x) + world.sky.sun.x
		suny := float32(world.camera.y) + world.sky.sun.y
		p := (sunx - suny - x + y)
		cell.biom.shine += 25 * exp(-p*p)
	}
}

func (cell *Cell) draw(camera *Camera, x, y float32) {
	if cell.dirt.height > 0 {
		cell.dirt.draw(camera, x, y)
	}
	cell.biom.draw(camera, x, y)
}
