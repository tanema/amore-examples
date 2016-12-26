package game

import (
	"github.com/tanema/amore-examples/platformer/ump"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"
)

type Player struct {
	x, y, width, height float32
	color               *gfx.Color
	body                *ump.Body
}

func NewPlayer(x, y, width, height float32, color *gfx.Color) *Player {
	newPlayer := &Player{
		x:      x,
		y:      y,
		width:  width,
		height: height,
		color:  color,
	}
	newPlayer.body = world.Add("player", x, y, width, height)
	return newPlayer
}

func (player *Player) Update(dt float32) {
	if keyboard.IsDown(keyboard.KeyUp) {
		player.y -= 1
	}
	if keyboard.IsDown(keyboard.KeyDown) {
		player.y += 1
	}
	if keyboard.IsDown(keyboard.KeyLeft) {
		player.x -= 1
	}
	if keyboard.IsDown(keyboard.KeyRight) {
		player.x += 1
	}
	player.x, player.y, _ = player.body.Move(player.x, player.y)
	camera.LookAt(player.x, player.y)
}

func (player *Player) Draw() {
	gfx.SetColorC(player.color)
	gfx.Rect(gfx.FILL, player.x, player.y, player.width, player.height)
}
