package game

import (
	"github.com/tanema/amore"
	"github.com/tanema/amore/keyboard"
)

// World encapsulates the whole environment
type World struct {
	size      int
	terrain   [][]*Cell
	camera    *Camera
	timeOfDay float32
	sin       float32
	sky       *Sky
	player    *Voxel
}

const (
	worldSaturation  float32 = 0.99
	baseShine        float32 = 0.4
	playerShineRange float32 = 25
	sunSize          float32 = 60
	moonSize         float32 = 30
)

// NewWorld generates a new world to render
func NewWorld(worldSize, visible, iterations int, smooth bool) *World {
	return &World{
		size:    worldSize,
		terrain: generateTerrain(worldSize, iterations, smooth),
		camera:  newCamera(visible),
		sky:     newSky(),
		player:  newVoxel(0, 0, 0, 1, 1, 0, worldSaturation, 0.5, true),
	}
}

func (world *World) getCell(x, y int) *Cell {
	// if x >= len(world.terrain) || x < 0 || y >= len(world.terrain[x]) || y < 0 {
	//	return nil
	// }
	i := (x + world.size) % world.size
	j := (y + world.size) % world.size
	return world.terrain[i][j]
}

// Update updates  a step in the world
func (world *World) Update(dt float32) {
	world.camera.update()
	world.updateInput()
	world.timeOfDay += dt / 10
	world.sin = sin(world.timeOfDay)
	world.sky.update(world)
	world.camera.forVisible(world, func(cell *Cell, x, y, distX, distY float32) {
		cell.update(world, x, y, distX, distY)
	})
}

// Draw draws one frame
func (world *World) Draw() {
	world.sky.draw(world)
	world.camera.forVisible(world, func(cell *Cell, x, y, distX, distY float32) {
		cell.draw(world.camera, x, y)
		if x == world.player.x && y == world.player.y {
			world.player.draw(world.camera, x, y)
		}
	})
}

func (world *World) updateInput() {
	if keyboard.IsDown(keyboard.KeyEscape) {
		amore.Quit()
	}

	if keyboard.IsDown(keyboard.KeyLeft) {
		world.player.x--
	} else if keyboard.IsDown(keyboard.KeyRight) {
		world.player.x++
	}
	if keyboard.IsDown(keyboard.KeyUp) {
		world.player.y--
	} else if keyboard.IsDown(keyboard.KeyDown) {
		world.player.y++
	}

	world.player.x = float32((int(world.player.x) + world.size) % world.size)
	world.player.y = float32((int(world.player.y) + world.size) % world.size)
	world.camera.lookAt(world.player.x, world.player.y)
	cell := world.getCell(int(world.player.x), int(world.player.y))

	if keyboard.IsDown(keyboard.KeySpace) {
		cell.setZ(cell.getZ() + 1)
	} else if keyboard.IsDown(keyboard.KeyC) {
		cell.setZ(cell.getZ() - 1)
	}

	world.player.z = cell.getZ()

	if keyboard.IsDown(keyboard.KeyV) {
		world.camera.visible++
	} else if keyboard.IsDown(keyboard.KeyB) {
		world.camera.visible--
	}
}
