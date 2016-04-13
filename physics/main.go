package main

import (
	"fmt"
	"math"
	"math/rand"

	b2d "github.com/neguse/go-box2d-lite/box2dlite"
	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/timer"
)

const timeStep = 1.0 / 60

var (
	gravity    = b2d.Vec2{0.0, -10.0}
	iterations = 10
	world      = b2d.NewWorld(gravity, iterations)
	title      = ""
)

func main() {
	Demo1()
	amore.Start(update, draw)
}

func update(deltaTime float32) {
	if keyboard.IsDown(keyboard.KeyEscape) {
		amore.Quit()
	}
	if keyboard.IsDown(keyboard.Key1) {
		Demo1()
	} else if keyboard.IsDown(keyboard.Key2) {
		Demo2()
	} else if keyboard.IsDown(keyboard.Key3) {
		Demo3()
	} else if keyboard.IsDown(keyboard.Key4) {
		Demo4()
	} else if keyboard.IsDown(keyboard.Key5) {
		Demo5()
	} else if keyboard.IsDown(keyboard.Key6) {
		Demo6()
	} else if keyboard.IsDown(keyboard.Key7) {
		Demo7()
	} else if keyboard.IsDown(keyboard.Key8) {
		Demo8()
	} else if keyboard.IsDown(keyboard.Key9) {
		Demo9()
	}
	world.Step(float64(deltaTime))
}

func renderBody(b *b2d.Body) {
	R := b2d.Mat22ByAngle(b.Rotation)
	x := b.Position
	h := b2d.MulSV(0.5, b.Width)

	o := b2d.Vec2{400, 400}
	S := b2d.Mat22{b2d.Vec2{20.0, 0.0}, b2d.Vec2{0.0, -20.0}}

	v1 := o.Add(S.MulV(x.Add(R.MulV(b2d.Vec2{-h.X, -h.Y}))))
	v2 := o.Add(S.MulV(x.Add(R.MulV(b2d.Vec2{h.X, -h.Y}))))
	v3 := o.Add(S.MulV(x.Add(R.MulV(b2d.Vec2{h.X, h.Y}))))
	v4 := o.Add(S.MulV(x.Add(R.MulV(b2d.Vec2{-h.X, h.Y}))))

	gfx.Line(
		float32(v1.X), float32(v1.Y),
		float32(v2.X), float32(v2.Y),
		float32(v3.X), float32(v3.Y),
		float32(v4.X), float32(v4.Y),
		float32(v1.X), float32(v1.Y))
}

func renderJoint(j *b2d.Joint) {
	b1 := j.Body1
	b2 := j.Body2

	R1 := b2d.Mat22ByAngle(b1.Rotation)
	R2 := b2d.Mat22ByAngle(b2.Rotation)

	x1 := b1.Position
	p1 := x1.Add(R1.MulV(j.LocalAnchor1))

	x2 := b2.Position
	p2 := x2.Add(R2.MulV(j.LocalAnchor2))

	o := b2d.Vec2{400, 400}
	S := b2d.Mat22{b2d.Vec2{20.0, 0.0}, b2d.Vec2{0.0, -20.0}}

	x1 = o.Add(S.MulV(x1))
	p1 = o.Add(S.MulV(p1))
	x2 = o.Add(S.MulV(x2))
	p2 = o.Add(S.MulV(p2))

	gfx.Line(float32(x1.X), float32(x1.Y), float32(p1.X), float32(p1.Y))
	gfx.Line(float32(x2.X), float32(x2.Y), float32(p2.X), float32(p2.Y))
}

func draw() {
	gfx.SetLineWidth(2)
	for _, b := range world.Bodies {
		renderBody(b)
	}
	for _, j := range world.Joints {
		renderJoint(j)
	}
	gfx.Print(fmt.Sprintf("fps: %v", timer.GetFPS()))
	gfx.Print("Press numbers 1 -> 9 to see different demos", 0, 15)
	gfx.Print(title, 100, 100)
}

// Single box
func Demo1() {
	title = "Single Box"
	world.Clear()

	var b1, b2 b2d.Body

	b1.Set(&b2d.Vec2{100.0, 20.0}, math.MaxFloat64)
	b1.Position = b2d.Vec2{0.0, -0.5 * b1.Width.Y}
	world.AddBody(&b1)

	b2.Set(&b2d.Vec2{1.0, 1.0}, 200.0)
	b2.Position = b2d.Vec2{0.0, 4.0}
	world.AddBody(&b2)
}

// A simple pendulum
func Demo2() {
	title = "Single Pendulum"
	world.Clear()

	var b2, b1 b2d.Body
	var j b2d.Joint

	b1.Set(&b2d.Vec2{100.0, 20.0}, math.MaxFloat64)
	b1.Friction = 0.2
	b1.Position = b2d.Vec2{0.0, -0.5 * b1.Width.Y}
	b1.Rotation = 0.0
	world.AddBody(&b1)

	b2.Set(&b2d.Vec2{1.0, 1.0}, 100.0)
	b2.Friction = 0.2
	b2.Position = b2d.Vec2{9.0, 11.0}
	b2.Rotation = 0.0
	world.AddBody(&b2)

	j.Set(&b1, &b2, &b2d.Vec2{0.0, 11.0})
	world.AddJoint(&j)
}

// Varying friction coefficients
func Demo3() {
	title = "Varying friction coefficients"
	world.Clear()

	{
		var b b2d.Body
		b.Set(&b2d.Vec2{100.0, 20.0}, math.MaxFloat64)
		b.Position = b2d.Vec2{0.0, -0.5 * b.Width.Y}
		world.AddBody(&b)
	}

	{
		var b b2d.Body
		b.Set(&b2d.Vec2{13.0, 0.25}, math.MaxFloat64)
		b.Position = b2d.Vec2{-2.0, 11.0}
		b.Rotation = -0.25
		world.AddBody(&b)
	}

	{
		var b b2d.Body
		b.Set(&b2d.Vec2{0.25, 1.0}, math.MaxFloat64)
		b.Position = b2d.Vec2{5.25, 9.5}
		world.AddBody(&b)
	}

	{
		var b b2d.Body
		b.Set(&b2d.Vec2{13.0, 0.25}, math.MaxFloat64)
		b.Position = b2d.Vec2{2.0, 7.0}
		b.Rotation = 0.25
		world.AddBody(&b)
	}

	{
		var b b2d.Body
		b.Set(&b2d.Vec2{0.25, 1.0}, math.MaxFloat64)
		b.Position = b2d.Vec2{-5.25, 5.5}
		world.AddBody(&b)
	}

	frictions := []float64{0.75, 0.5, 0.35, 0.1, 0.0}
	for i := 0; i < 5; i++ {
		var b b2d.Body
		b.Set(&b2d.Vec2{0.5, 0.5}, 25.0)
		b.Friction = frictions[i]
		b.Position = b2d.Vec2{-7.5 + 2.0*float64(i), 14.0}
		world.AddBody(&b)
	}

}

// A vertical stack
func Demo4() {
	title = "A vertical stack"
	world.Clear()

	{
		var b b2d.Body
		b.Set(&b2d.Vec2{100.0, 20.0}, math.MaxFloat64)
		b.Friction = 0.2
		b.Position = b2d.Vec2{0.0, -0.5 * b.Width.Y}
		b.Rotation = 0.0
		world.AddBody(&b)
	}

	for i := 0; i < 10; i++ {
		var b b2d.Body
		b.Set(&b2d.Vec2{1.0, 1.0}, 1.0)
		b.Friction = 0.2
		x := rand.Float64()*0.2 - 0.1
		b.Position = b2d.Vec2{x, 0.51 + 1.05*float64(i)}
		world.AddBody(&b)
	}

}

// A pyramid
func Demo5() {
	title = "A pyramid"
	world.Clear()

	{
		var b b2d.Body
		b.Set(&b2d.Vec2{100.0, 20.0}, math.MaxFloat64)
		b.Friction = 0.2
		b.Position = b2d.Vec2{0.0, -0.5 * b.Width.Y}
		b.Rotation = 0.0
		world.AddBody(&b)
	}

	x := b2d.Vec2{-6.0, 0.75}

	for i := 0; i < 12; i++ {
		y := x
		for j := i; j < 12; j++ {
			var b b2d.Body
			b.Set(&b2d.Vec2{1.0, 1.0}, 10.0)
			b.Friction = 0.2
			b.Position = y
			world.AddBody(&b)

			y = y.Add(b2d.Vec2{1.125, 0.0})
		}

		x = x.Add(b2d.Vec2{0.5625, 2.0})
	}
}

// A teeter
func Demo6() {
	title = "A teeter"
	world.Clear()

	var b1, b2, b3, b4, b5 b2d.Body
	b1.Set(&b2d.Vec2{100.0, 20.0}, math.MaxFloat64)
	b1.Position = b2d.Vec2{0.0, -0.5 * b1.Width.Y}
	world.AddBody(&b1)

	b2.Set(&b2d.Vec2{12.0, 0.25}, 100)
	b2.Position = b2d.Vec2{0.0, 1.0}
	world.AddBody(&b2)

	b3.Set(&b2d.Vec2{0.5, 0.5}, 25.0)
	b3.Position = b2d.Vec2{-5.0, 2.0}
	world.AddBody(&b3)

	b4.Set(&b2d.Vec2{0.5, 0.5}, 25.0)
	b4.Position = b2d.Vec2{-5.5, 2.0}
	world.AddBody(&b4)

	b5.Set(&b2d.Vec2{1.0, 1.0}, 100)
	b5.Position = b2d.Vec2{5.5, 15.0}
	world.AddBody(&b5)

	{
		var j b2d.Joint
		j.Set(&b1, &b2, &b2d.Vec2{0.0, 1.0})
		world.AddJoint(&j)
	}

}

// A suspension bridge
func Demo7() {
	title = "A suspension bridge"
	world.Clear()

	var ba []*b2d.Body

	{
		var b b2d.Body
		b.Set(&b2d.Vec2{100.0, 20.0}, math.MaxFloat64)
		b.Friction = 0.2
		b.Position = b2d.Vec2{0.0, -0.5 * b.Width.Y}
		b.Rotation = 0.0
		world.AddBody(&b)
		ba = append(ba, &b)
	}

	const numPlunks = 15
	const mass = 50.0

	for i := 0; i < numPlunks; i++ {
		var b b2d.Body
		b.Set(&b2d.Vec2{1.0, 0.25}, mass)
		b.Friction = 0.2
		b.Position = b2d.Vec2{-8.5 + 1.25*float64(i), 5.0}
		world.AddBody(&b)
		ba = append(ba, &b)
	}

	// Tuning
	const frequencyHz = 2.0
	const dampingRatio = 0.7

	// frequency in radians
	const omega = 2.0 * math.Pi * frequencyHz

	// damping coefficient
	const d = 2.0 * mass * dampingRatio * omega

	// spring stifness
	const k = mass * omega * omega

	// magic formulas
	const softness = 1.0 / (d + timeStep*k)
	const biasFactor = timeStep * k / (d + timeStep*k)

	for i := 0; i <= numPlunks; i++ {
		var j b2d.Joint
		j.Set(ba[i], ba[(i+1)%(numPlunks+1)], &b2d.Vec2{-9.125 + 1.25*float64(i), 5.0})
		j.Softness = softness
		j.BiasFactor = biasFactor
		world.AddJoint(&j)
	}

}

// Dominos
func Demo8() {
	title = "Dominos"
	world.Clear()

	var b1 b2d.Body
	b1.Set(&b2d.Vec2{100.0, 20.0}, math.MaxFloat64)
	b1.Position = b2d.Vec2{0.0, -0.5 * b1.Width.Y}
	world.AddBody(&b1)

	{
		var b b2d.Body
		b.Set(&b2d.Vec2{12.0, 0.5}, math.MaxFloat64)
		b.Position = b2d.Vec2{-1.5, 10.0}
		world.AddBody(&b)
	}

	for i := 0; i < 10; i++ {
		var b b2d.Body
		b.Set(&b2d.Vec2{0.2, 2.0}, 10.0)
		b.Position = b2d.Vec2{-6.0 + 1.0*float64(i), 11.25}
		b.Friction = 0.1
		world.AddBody(&b)
	}

	{
		var b b2d.Body
		b.Set(&b2d.Vec2{14.0, 0.5}, math.MaxFloat64)
		b.Position = b2d.Vec2{1.0, 6.0}
		b.Rotation = 0.3
		world.AddBody(&b)
	}

	var b2 b2d.Body
	b2.Set(&b2d.Vec2{0.5, 3.0}, math.MaxFloat64)
	b2.Position = b2d.Vec2{-7.0, 4.0}
	world.AddBody(&b2)

	var b3 b2d.Body
	b3.Set(&b2d.Vec2{12.0, 0.25}, 20.0)
	b3.Position = b2d.Vec2{-0.9, 1.0}
	world.AddBody(&b3)

	{
		var j b2d.Joint
		j.Set(&b1, &b3, &b2d.Vec2{-2.0, 1.0})
		world.AddJoint(&j)
	}

	var b4 b2d.Body
	b4.Set(&b2d.Vec2{0.5, 0.5}, 10.0)
	b4.Position = b2d.Vec2{-10.0, 15.0}
	world.AddBody(&b4)

	{
		var j b2d.Joint
		j.Set(&b2, &b4, &b2d.Vec2{-7.0, 15.0})
		world.AddJoint(&j)
	}

	var b5 b2d.Body
	b5.Set(&b2d.Vec2{2.0, 2.0}, 20.0)
	b5.Position = b2d.Vec2{6.0, 2.5}
	b5.Friction = 0.1
	world.AddBody(&b5)

	{
		var j b2d.Joint
		j.Set(&b1, &b5, &b2d.Vec2{6.0, 2.6})
		world.AddJoint(&j)
	}

	var b6 b2d.Body
	b6.Set(&b2d.Vec2{2.0, 0.2}, 10.0)
	b6.Position = b2d.Vec2{6.0, 3.6}
	world.AddBody(&b6)

	{
		var j b2d.Joint
		j.Set(&b5, &b6, &b2d.Vec2{7.0, 3.5})
		world.AddJoint(&j)
	}

}

// A multi-pendulum
func Demo9() {
	title = "A multi-pendulum"
	world.Clear()

	var b1 *b2d.Body

	{
		var b b2d.Body
		b.Set(&b2d.Vec2{100.0, 20.0}, math.MaxFloat64)
		b.Position = b2d.Vec2{0.0, -0.5 * b.Width.Y}
		world.AddBody(&b)
		b1 = &b
	}

	const mass = 10.0

	// Tuning
	const frequencyHz = 4.0
	const dampingRatio = 0.7

	// frequency in radians
	const omega = 2.0 * math.Pi * frequencyHz

	// damping coefficient
	const d = 2.0 * mass * dampingRatio * omega

	// spring stiffness
	const k = mass * omega * omega

	// magic formulas
	const softness = 1.0 / (d + timeStep*k)
	const biasFactor = timeStep * k / (d + timeStep*k)

	const y = 12.0

	for i := 0; i < 15; i++ {
		x := b2d.Vec2{0.5 + float64(i), y}

		var b b2d.Body
		b.Set(&b2d.Vec2{0.75, 0.25}, mass)
		b.Friction = 0.2
		b.Position = x
		b.Rotation = 0.0
		world.AddBody(&b)

		var j b2d.Joint
		j.Set(b1, &b, &b2d.Vec2{float64(i), y})
		j.Softness = softness
		j.BiasFactor = biasFactor
		world.AddJoint(&j)

		b1 = &b
	}

}
