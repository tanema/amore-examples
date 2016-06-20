package game

import (
	"github.com/tanema/amore-examples/platformer/lib/bump"
	//"github.com/tanema/amore-examples/platformer/lib/gamera"
)

const (
	updateRadius float32 = 100 // how "far away from the camera" things stop being updated
)

var (
	width  float32 = 4000
	height float32 = 2000
	//camera *gamera.Camera
	player *Player
	block  *Block
	world  *bump.World
)

func New() {
	world = bump.NewWorld(64)
	player = newPlayer(world, 60, 60)
	block = newBlock(world, 0, 100, 100, 50)
	//camera = gamera.New(0, 0, width, height)
}

func Update(dt float32) {
	l, t, w, h := float32(0), float32(0), float32(800), float32(600) //camera.GetVisible()
	l, t, w, h = l-updateRadius, t-updateRadius, w+updateRadius*2, h+updateRadius*2
	for _, item := range world.QueryRect(l, t, w, h) {
		item.Entity.Update(dt)
	}
	//camera.SetPosition(player.l, player.t)
	//camera.Update(dt)
}

func Draw() {
	l, t, w, h := float32(0), float32(0), float32(800), float32(600)
	//camera.Draw(func(l, t, w, h float32) {
	for _, item := range world.QueryRect(l, t, w, h) {
		item.Entity.Draw()
	}
	//})
}
