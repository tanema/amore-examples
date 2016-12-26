package game

import (
	"fmt"

	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"

	"github.com/tanema/amore-examples/platformer/lense"
	"github.com/tanema/amore-examples/platformer/ump"
)

const (
	updateRadius float32 = 100 // how "far away from the camera" things stop being updated
)

var (
	width  float32 = 4000
	height float32 = 2000
	camera *lense.Camera
	scale  = float32(1)
	rot    = float32(0)
	world  = ump.NewWorld(64)

	blocks = []*Block{}
	player *Player
)

func New() {
	camera = lense.New()
	player = NewPlayer(200, 200, 50, 50, gfx.NewColor(255, 0, 0, 255))
	for i := 0; i <= 10; i++ {
		blocks = append(blocks, NewBlock(
			randRange(0, 800), randRange(0, 600), randRange(50, 200), randRange(50, 200),
			gfx.NewColor(randRange(0, 255), randRange(0, 255), randRange(0, 255), 255),
		))
	}
}

func Update(dt float32) {
	if keyboard.IsDown(keyboard.KeyEscape) {
		amore.Quit()
	}
	if keyboard.IsDown(keyboard.KeyE) {
		println("shake")
		camera.Shake(1)
	}
	if keyboard.IsDown(keyboard.KeyW) {
		scale += 0.01
	}
	if keyboard.IsDown(keyboard.KeyS) {
		scale -= 0.01
	}
	if keyboard.IsDown(keyboard.KeyD) {
		rot += 0.01
	}
	if keyboard.IsDown(keyboard.KeyA) {
		rot -= 0.01
	}

	//l, t, w, h := camera.GetVisible()
	//for _, item := range world.QueryRect(l-updateRadius, t-updateRadius, w+updateRadius*2, h+updateRadius*2) {
	//item.Entity.Update(dt)
	//}
	player.Update(dt)
	fmt.Println(scale)
	camera.ZoomTo(scale)
	camera.RotateTo(rot)
	camera.Update(dt)
}

func Draw() {
	camera.Draw(func(l, t, w, h float32) {
		world.DrawDebug(0, 0, 800, 600)
		for _, block := range blocks {
			block.Draw()
		}
		player.Draw()
	})
	//for _, item := range world.QueryRect(0, 0, 800, 600) {
	//item.Entity.Draw()
	//item.DrawDebug()
	//}
}
