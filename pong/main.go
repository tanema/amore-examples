// Simple pong clone without more complex paddle physics. You can tune the ball
// speed and paddle speed for better playability
package main

import (
	"fmt"

	"github.com/tanema/amore"
	"github.com/tanema/amore/audio"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/window"
)

const (
	paddleWidth  = 10
	paddleHeigth = 100
	paddleSpeed  = 200
	ballSpeed    = 200
)

var (
	gameOver              = false
	enemyPosition float32 = 0
	screenWidth   float32
	screenHeight  float32

	player = &Paddle{
		x:      800 - paddleWidth,
		y:      300 - (paddleHeigth / 2),
		width:  paddleWidth,
		height: paddleHeigth,
		speed:  paddleSpeed,
		color:  gfx.NewColor(0, 255, 0, 255),
	}

	enemy = &Paddle{
		y:      300 - (paddleHeigth / 2),
		width:  paddleWidth,
		height: paddleHeigth,
		speed:  paddleSpeed,
		color:  gfx.NewColor(0, 0, 255, 255),
	}

	ball = &Ball{
		x:      400,
		y:      300,
		vx:     1.5,
		vy:     1.5,
		radius: 10,
		speed:  ballSpeed,
		color:  gfx.NewColor(255, 0, 0, 255),
	}

	scoreLabel *gfx.Text
	blip, _    = audio.NewSource("audio/blip.wav", true)
	bomb, _    = audio.NewSource("../test-all/assets/audio/bomb.wav", true)
)

func main() {
	window.SetMouseVisible(false)
	amore.OnLoad = load
	amore.Start(update, draw)
}

func load() {
	screenWidth = gfx.GetWidth()
	screenHeight = gfx.GetHeight()
	scoreLabel, _ = gfx.NewText(
		gfx.GetFont(),
		fmt.Sprintf("%v : %v", enemy.score, player.score),
	)
}

func update(dt float32) {
	if keyboard.IsDown(keyboard.KeyEscape) {
		amore.Quit()
	}
	if keyboard.IsDown(keyboard.KeyReturn) {
		reset()
	}

	if player.score >= 10 || enemy.score >= 10 {
		gameOver = true
	}

	if gameOver {
		return
	}

	ball.Update(dt)

	// enemy movements
	enemyY := enemy.y + (enemy.height / 3)
	if ball.y < enemyY {
		enemy.y -= enemy.speed * dt
	} else if ball.y > enemyY {
		enemy.y += enemy.speed * dt
	}

	//if keyboard.IsDown(keyboard.KeyW) {
	//enemy.y -= enemy.speed * dt
	//} else if keyboard.IsDown(keyboard.KeyS) {
	//enemy.y += enemy.speed * dt
	//}

	// player movements
	if keyboard.IsDown(keyboard.KeyUp) {
		player.y -= player.speed * dt
	} else if keyboard.IsDown(keyboard.KeyDown) {
		player.y += player.speed * dt
	}

	scoreLabel.Set(fmt.Sprintf("%v : %v", enemy.score, player.score))
}

func reset() {
	ball.Reset()

	player.y = 300 - (paddleHeigth / 2)
	player.score = 0

	enemy.y = 300 - (paddleHeigth / 2)
	enemy.score = 0

	gameOver = false
}

func draw() {
	halfWidth := scoreLabel.GetWidth() / 2
	gfx.SetColor(255, 255, 255, 255)
	scoreLabel.Draw((screenWidth/2)-halfWidth, 0)
	gfx.SetLineWidth(3)
	gfx.Line(screenWidth/2+1, 25, screenWidth/2+1, screenHeight)

	ball.Draw()
	enemy.Draw()
	player.Draw()

	if gameOver {
		gameOverString := "Game Over"
		wonlost := "You Won"
		if enemy.score > player.score {
			wonlost = "You Lost"
		}
		pressEnter := "Press Enter To Restart"
		enterWidth := gfx.GetFont().GetWidth(pressEnter) + 20

		gfx.SetColor(0, 0, 0, 255)
		leftAlign := gfx.GetWidth()/2 - (enterWidth / 2)

		gfx.Rect(gfx.FILL, leftAlign, 200, enterWidth, 75)
		if enemy.score > player.score {
			gfx.SetColorC(enemy.color)
		} else {
			gfx.SetColorC(player.color)
		}
		gfx.Rect(gfx.LINE, leftAlign, 200, enterWidth, 75)

		gfx.Print(gameOverString, leftAlign+10, 215)
		gfx.Print(wonlost, leftAlign+10, 230)
		gfx.Print(pressEnter, leftAlign+10, 245)
	}
}

type Ball struct {
	x      float32
	y      float32
	vx     float32
	vy     float32
	radius float32
	speed  float32
	color  *gfx.Color
}

func (ball *Ball) Reset() {
	ball.x = 400
	ball.y = 300
}

func (ball *Ball) Update(dt float32) {
	ball.x += ball.speed * ball.vx * dt
	ball.y += ball.speed * ball.vy * dt

	// enemy's wall
	if ball.x <= ball.radius {
		// if the ball is past the wall/paddle lets put it before the wall so
		// that we wont have multiple collisions
		if ball.x < ball.radius {
			ball.x = ball.radius
		}

		if ball.y >= enemy.y && ball.y <= enemy.y+enemy.height {
			blip.Play()
			ball.vx = ball.vx * -1
		} else {
			ball.Reset()
			bomb.Play()
			player.score += 1
		}
	}

	if ball.x >= screenWidth-ball.radius { //players wall
		// if the ball is past the wall/paddle lets put it before the wall so
		// that we wont have multiple collisions
		if ball.x > screenWidth-ball.radius {
			ball.x = screenWidth - ball.radius
		}

		if ball.y >= player.y && ball.y <= player.y+player.height {
			blip.Play()
			ball.vx = ball.vx * -1
		} else {
			ball.Reset()
			bomb.Play()
			enemy.score += 1
		}
	}

	if ball.y <= ball.radius || ball.y >= screenHeight-ball.radius {
		// if the ball is past the wall lets put it before the wall so
		// that we wont have multiple collisions
		if ball.y < ball.radius {
			ball.y = ball.radius
		}
		if ball.y > screenHeight-ball.radius {
			ball.y = screenHeight - ball.radius
		}

		blip.Play()
		ball.vy = ball.vy * -1
	}
}

func (ball *Ball) Draw() {
	gfx.SetColorC(ball.color)
	gfx.Circle(gfx.FILL, ball.x, ball.y, ball.radius)
}

type Paddle struct {
	score  int
	x      float32
	y      float32
	vy     float32
	width  float32
	height float32
	speed  float32
	color  *gfx.Color
}

func (paddle *Paddle) Draw() {
	gfx.SetColorC(paddle.color)
	gfx.Rect(gfx.FILL, paddle.x, paddle.y, paddle.width, paddle.height)
}
