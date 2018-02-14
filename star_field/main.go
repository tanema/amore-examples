package main

import (
	"math/rand"
	"time"

	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
)

const (
	starSpawnDepth float32 = 20
	starSpeed      float32 = 0.05
	screenWidth    float32 = 800 // Default screen size
	screenHeight   float32 = 600
)

type star struct {
	x  float32
	y  float32
	z  float32
	px float32
	py float32
	cx float32
	cy float32
}

func randMax(max float32) float32 {
	rand.Seed(time.Now().UTC().UnixNano())
	return (rand.Float32() * (max + 1))
}

func (s *star) reset() {
	s.x = (randMax(screenWidth) - (screenWidth * 0.5)) * starSpawnDepth
	s.y = (randMax(screenWidth) - (screenHeight * 0.5)) * starSpawnDepth
	s.z = starSpawnDepth
	s.px = 0
	s.py = 0
}

func (s *star) update(dt float32) {
	s.px = s.x / s.z
	s.py = s.y / s.z
	s.z -= starSpeed
	s.cx = screenWidth * 0.5
	s.cy = screenHeight * 0.5
	if s.z < 0 || s.px > screenWidth || s.py > screenHeight {
		s.reset()
	}
}

func (s *star) draw() {
	if s.px == 0 {
		return
	}
	x := s.x / s.z
	y := s.y / s.z
	gfx.SetLineWidth((1.0/s.z + 1) * 2)
	gfx.PolyLine([]float32{x + s.cx, y + s.cy, s.px + s.cx, s.py + s.cy})
}

func main() {
	units := int(300)
	stars := make([]*star, units)
	for i := 0; i < units; i++ {
		stars[i] = &star{}
		stars[i].reset()
	}

	amore.Start(func(dt float32) {
		for _, s := range stars {
			s.update(dt)
		}
	}, func() {
		gfx.SetColor(255, 255, 255, 255)
		for _, s := range stars {
			s.draw()
		}
	})
}
