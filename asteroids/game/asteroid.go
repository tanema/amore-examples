package game

const (
	asteroidSpeed = 100
	asteroidSpin  = 2
)

var asteroidPoints = []float32{
	-10, 0,
	-5, 7,
	-3, 4,
	1, 10,
	5, 4,
	10, 0,
	5, -6,
	2, -10,
	-4, -10,
	-4, -5,
	-10, 0,
}

type Asteroid struct {
	*Sprite
	parent bool
}

func newAsteroid() *Asteroid {
	new_asteroid := &Asteroid{
		parent: true,
	}

	new_asteroid.Sprite = NewSprite(new_asteroid, "asteroid",
		randMax(screenWidth),
		randMax(screenHeight),
		randRange(3, 8),
		asteroidPoints, true)
	new_asteroid.vx = randLimits(asteroidSpeed)
	new_asteroid.vy = randLimits(asteroidSpeed)
	new_asteroid.vrot = randLimits(asteroidSpin)

	return new_asteroid
}

func (asteroid *Asteroid) Update(dt float32) {
	if collisions := asteroid.UpdateMovement(dt); len(collisions) > 0 {
		for _, c := range collisions {
			if collisions[0].Name == "ship" {
				c.Collidable.Destroy(false)
			} else if collisions[0].Name == "bullet" {
				score++
				asteroid.Destroy(false)
				c.Collidable.Destroy(false)
			}
		}
	}
}

func (asteroid *Asteroid) Destroy(force bool) {
	removeObject(asteroid)
	asteroid.Sprite.Destroy()
	if !force {
		if asteroid.parent {
			for i := 0; i < 3; i++ {
				a := &Asteroid{parent: false}
				a.Sprite = NewSprite(a, "asteroid",
					asteroid.x,
					asteroid.y,
					randRange(1, 3),
					asteroidPoints, true)
				a.vx = randLimits(asteroidSpeed)
				a.vy = randLimits(asteroidSpeed)
				a.vrot = randLimits(asteroidSpin)
				addObject(a)
			}
		} else {
			newExplosion(asteroid.Sprite)
		}
		bomb.Play()
	}
}
