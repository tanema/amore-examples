package game

import (
	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"

	"github.com/tanema/amore-examples/platformer/lib/bump"
)

const (
	updateRadius float32 = 100 // how "far away from the camera" things stop being updated
)

var (
	width  float32 = 4000
	height float32 = 2000
	camera *Camera
	x      = float32(200)
	y      = float32(200)
	scale  = float32(1)
	rot    = float32(0)
	world  *bump.World

	blocks = []*Block{}
	player *Block
)

func New() {
	world = bump.NewWorld(64)
	camera = NewCamera()

	player = NewBlock(x, y, 50, 50, gfx.NewColor(255, 0, 0, 255))
	world.Add(player, "player", x, y, 50, 50, map[string]string{})

	for i := 0; i <= 1; i++ {
		bx, by, bw, bh := randRange(0, 800), randRange(0, 600), randRange(50, 200), randRange(50, 200)

		blocks = append(blocks, NewBlock(
			bx, by, bw, bh,
			gfx.NewColor(randRange(0, 255), randRange(0, 255), randRange(0, 255), 255),
		))
		world.Add(blocks[len(blocks)-1], "block", bx, by, bw, bh, map[string]string{})
	}
}

func Update(dt float32) {
	if keyboard.IsDown(keyboard.KeyEscape) {
		amore.Quit()
	}
	if keyboard.IsDown(keyboard.KeyUp) {
		y -= 1
	}
	if keyboard.IsDown(keyboard.KeyDown) {
		y += 1
	}
	if keyboard.IsDown(keyboard.KeyLeft) {
		x -= 1
	}
	if keyboard.IsDown(keyboard.KeyRight) {
		x += 1
	}
	if keyboard.IsDown(keyboard.KeyE) {
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

	l, t, w, h := camera.GetVisible()
	for _, item := range world.QueryRect(l-updateRadius, t-updateRadius, w+updateRadius*2, h+updateRadius*2) {
		item.Entity.Update(dt)
	}
	x, y, _ = world.Move(player, x, y)
	player.x, player.y = x, y
	camera.LookAt(x, y)
	camera.ZoomTo(scale)
	camera.RotateTo(rot)
	camera.Update(dt)
}

func Draw() {
	camera.Draw(func(l, t, w, h float32) {
		world.DrawDebug(l, t, w, h)
		for _, item := range world.QueryRect(l, t, w, h) {
			item.Entity.Draw()
			item.DrawDebug()
		}
	})
}
