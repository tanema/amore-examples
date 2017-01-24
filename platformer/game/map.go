package game

import (
	"github.com/tanema/amore-examples/platformer/lense"
	"github.com/tanema/amore-examples/platformer/ump"
)

type gameObject interface {
	Update(dt float32)
	tag() string
	destroy()
	push(strength float32)
	damage(intensity float32)
	Draw(bool)
}

type Map struct {
	width        float32
	height       float32
	updateRadius float32
	Player       *Player
	objects      map[uint32]gameObject
	debug        bool
	camera       *lense.Camera
	world        *ump.World
}

func NewMap(width, height float32, camera *lense.Camera) *Map {
	gameMap := &Map{
		width:        width,
		height:       height,
		updateRadius: 100,
		camera:       camera,
	}
	gameMap.Reset()
	return gameMap
}

func (m *Map) Reset() {
	m.objects = map[uint32]gameObject{}
	m.world = ump.NewWorld(64)
	m.Player = newPlayer(m, 60, 60)

	// walls & ceiling
	newBlock(m, 0, 0, m.width, 32, true)
	newBlock(m, 0, 32, 32, m.height-64, true)
	newBlock(m, m.width-32, 32, 32, m.height-64, true)

	// tiled floor
	tilesOnFloor := float32(40)
	for i := float32(0); i < tilesOnFloor; i++ {
		newBlock(m, i*m.width/tilesOnFloor, m.height-32, m.width/tilesOnFloor, 32, true)
	}

	// groups of blocks
	for i := 1; i < 60; i++ {
		w := randRange(100, 400)
		h := randRange(100, 400)
		area := w * h
		l := randRange(100, m.width-w-200)
		t := randRange(100, m.height-h-100)
		indestructible := randRange(0, 1) < 0.75

		for j := 1; j < int(floor(area/7000)); j++ {
			newBlock(m,
				randRange(l, l+w),
				randRange(t, t+h),
				randRange(32, 100),
				randRange(32, 100),
				indestructible)
		}
	}

	for i := 0; i < 10; i++ {
		newGuardian(m,
			randRange(100, m.width-200),
			randRange(100, m.height-150))
	}
}

func (m *Map) ToggleDebug() {
	m.debug = !m.debug
}

func (m *Map) Update(dt, l, t, w, h float32) {
	if m.Player.isDead {
		m.Player.deadCounter = m.Player.deadCounter + dt
		if m.Player.deadCounter >= deadDuration {
			m.Reset()
		}
	}

	l, t, w, h = l-m.updateRadius, t-m.updateRadius, w+m.updateRadius*2, h+m.updateRadius*2
	for _, item := range m.world.QueryRect(l, t, w, h) {
		object, ok := m.objects[item.ID]
		if ok {
			object.Update(dt)
		}
	}
}

func (m *Map) Draw(l, t, w, h float32) {
	if m.debug {
		m.world.DrawDebug(l, t, w, h)
	}
	for _, item := range m.world.QueryRect(l, t, w, h) {
		object, ok := m.objects[item.ID]
		if ok {
			object.Draw(m.debug)
		}
	}
}

func (m *Map) Get(item *ump.Body) gameObject {
	return m.objects[item.ID]
}
