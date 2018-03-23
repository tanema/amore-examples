package main

import (
	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"

	"github.com/tanema/amore-examples/2d3d/render"
)

var (
	triforce *render.Fake3D
	rotation float32
)

const rotationSpeed float32 = 2

func main() {
	triforce, _ = render.New("images/triforce.png", 16, 16)
	amore.Start(update, draw)
}

func update(dt float32) {
	rotation += rotationSpeed * dt
}

func draw() {
	gfx.Scale(10)
	triforce.Draw(30, 30, rotation)
}
