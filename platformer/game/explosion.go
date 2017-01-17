package game

var (
	explosionWidth        float32 = 150
	explosionHeight       float32 = explosionWidth
	explosionMaxPushSpeed float32 = 300
)

func newExplosion(grenade *Grenade) {
	x, y := grenade.GetCenter()
	l, t, w, h := x-explosionWidth/2, y-explosionHeight/2, explosionWidth, explosionHeight
	gameMap := grenade.parent.gameMap
	world := gameMap.world
	gameMap.camera.Shake(6)

	for _, item := range world.QueryRect(l, t, w, h) {
		object := gameMap.Get(item)
		tag := object.tag()
		if tag == "player" || tag == "guardian" || tag == "block" {
			object.damage(0.7)
		}
	}

	radius := float32(50)
	for _, item := range world.QueryRect(l-radius, t-radius, w+radius+radius, h+radius+radius) {
		object := gameMap.Get(item)
		tag := object.tag()
		if tag == "player" || tag == "grenade" || tag == "debris" || tag == "puff" {
			object.push(300)
		}
	}

	for i := float32(0); i < randRange(15, 30); i++ {
		newPuff(
			grenade.parent.gameMap,
			randRange(l, l+w), randRange(t, t+h),
			0, -10, 2, 10,
		)
	}
}
