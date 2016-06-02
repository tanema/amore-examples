package game

import (
	"github.com/tanema/amore-examples/platformer/lib/bump"
	"github.com/tanema/amore-examples/platformer/lib/gamera"
)

type Map struct {
	player *Player
	width  float32
	height float32
	camera *gamera.Camera
}

func newMap(width, height float32, camera *gamera.Camera) *Map {
	return &Map{
		width:  width,
		height: height,
		camera: camera,
	}
}

func (m *Map) reset() {
	m.world = bump.NewWorld(64)
	m.player = newPlayer(m, m.world, 60, 60)

	// walls & ceiling
	newBlock(m.world, 0, 0, m.width, 32, true)
	newBlock(m.world, 0, 32, 32, m.height-64, true)
	newBlock(m.world, m.width-32, 32, 32, m.height-64, true)

	// tiled floor
	tilesOnFloor := 40
	for i := 0; i <= tilesOnFloor-1; i++ {
		newBlock(m.world, i*m.width/tilesOnFloor, m.height-32, m.width/tilesOnFloor, 32, true)
	}

	// groups of blocks
	for i := 1; i <= 60; i++ {
		w := randRange(100, 400)
		h := randRange(100, 400)
		area := w * h
		l := randRange(100, m.width-w-200)
		t := randRange(100, m.height-h-100)

		for i = 1; i <= math.floor(area/7000); i++ {
			newBlock(m.world, randRange(l, l+w), randRange(t, t+h), randRange(32, 100), randRange(32, 100), randRange() > 0.75)
		}
	}

	for i := 1; i <= 10; i++ {
		newGuardian(m.world, m.player, m.camera, randRange(100, m.width-200), randRange(100, m.height-150))
	}
}

func (m *Map) update(dt, l, t, w, h float32) {
	for _, item := range m.world.QueryRect(l, t, w, h) {
		item.update(dt)
	}
}

func (m *Map) draw(debug bool, t, w, h float32) {
	if debug {
		m.world.drawDebug(l, t, w, h)
	}
	for _, item := range m.world.QueryRect(l, t, w, h) {
		item.draw(debug)
	}
}
