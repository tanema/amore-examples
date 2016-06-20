package game

import (
	"github.com/tanema/amore-examples/platformer/lib/bump"
)

type Player struct {
	*Entity
}

func newPlayer(world *bump.World, x, y float32) *Player {
	player := &Player{
		Entity: newEntity(world, x, y, 32, 64),
	}
	player.Body = world.Add(player, "player", x, y, 32, 64, map[string]string{})
	return player
}

func (player *Player) Send(event string, args ...interface{}) {
}

func (player *Player) Update(dt float32) {
}
