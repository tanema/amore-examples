package game

import (
	"fmt"

	"github.com/tanema/amore-examples/asteroids/game/phys"

	"github.com/tanema/amore/audio"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"
)

var debug bool

type GameObject interface {
	Update(dt float32)
	Draw()
	Destroy(force bool)
}

const (
	cellSize float32 = 60
)

var (
	score        = 0
	world        *phys.World
	objects      []GameObject
	bomb, _      = audio.NewSource("../test-all/assets/audio/bomb.wav", true)
	lazer, _     = audio.NewSource("audio/lazer.wav", true)
	player       *Player
	gameOver     = false
	screenWidth  float32
	screenHeight float32
)

func New() {
	keyboard.OnKeyUp = keyup
	screenWidth = gfx.GetWidth()
	screenHeight = gfx.GetHeight()
	gameOver = false
	score = 0
	world = phys.NewWorld(screenWidth, screenHeight, cellSize)
	player = newPlayer()
	objects = []GameObject{
		player,
		newAsteroid(),
		newAsteroid(),
		newAsteroid(),
		newAsteroid(),
		newAsteroid(),
	}
}

func keyup(key keyboard.Key) {
	if key == keyboard.KeyTab {
		debug = !debug
	}
}

func Update(dt float32) {
	for _, object := range objects {
		object.Update(dt)
	}
	gameOver = gameOver || len(objects) == 1
}

func addObject(object GameObject) {
	objects = append(objects, object)
}

func removeObject(object GameObject) {
	for i, other := range objects {
		if object == other {
			objects = append(objects[:i], objects[i+1:]...)
			return
		}
	}
}

func Draw() {
	if debug {
		world.DrawGrid()
		gfx.Print(fmt.Sprintf("objects: %v", len(objects)), 0, 15)
	}

	for _, object := range objects {
		object.Draw()
	}
	gfx.Print(fmt.Sprintf("Score: %v", score), screenWidth-75, 0)
	if gameOver {
		if len(objects) == 1 && objects[0] == player {
			gfx.Print("Game Over. You Won.", screenWidth/2-175, screenHeight/2, 0, 2, 3)
		} else {
			gfx.Print("Game Over. You Lost.", screenWidth/2-175, screenHeight/2, 0, 2, 2)
		}
	}
}
