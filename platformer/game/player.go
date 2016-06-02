package game

type Player struct {
	*Entity
	health      float32
	deadCounter int
	game_map    *Map
}

const (
	player_width  float32 = 32
	player_height float32 = 64
)

func newPlayer(game_map, world, x, y float32) *Player {
	return &Player{
		Entity:   newEntity(world, "player", x, y, player_width, player_height, map[string]string{}),
		game_map: game_map,
		health:   1,
	}
}
