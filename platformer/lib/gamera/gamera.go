package gamera

import (
	"math"
	"math/rand"

	"github.com/tanema/amore/gfx"
)

const (
	maxShake        = float64(5)
	atenuationSpeed = float32(4)
)

type Camera struct {
	x, y           float32
	l, t, w, h     float32
	wl, wt, ww, wh float32
	w2, h2         float32
	sin, cos       float32
	scale          float32
	angle          float32
	shakeIntensity float32
}

func New(l, t, w, h float32) *Camera {
	sw, sh := gfx.GetWidth(), gfx.GetHeight()

	cam := &Camera{
		x:     0,
		y:     0,
		scale: 1,
		angle: 0,
		sin:   float32(math.Sin(0)),
		cos:   float32(math.Cos(0)),
		l:     0,
		t:     0,
		w:     sw,
		h:     sh,
		w2:    sw * 0.5,
		h2:    sh * 0.5,
	}

	cam.SetWorld(l, t, w, h)

	return cam
}

func (camera *Camera) getVisibleArea() (float32, float32) {
	sin, cos := float32(math.Abs(float64(camera.sin))), float32(math.Abs(float64(camera.cos)))
	w, h := camera.w/camera.scale, camera.h/camera.scale
	w, h = cos*w+sin*h, sin*w+cos*h
	return float32(math.Min(float64(w), float64(camera.ww))),
		float32(math.Min(float64(h), float64(camera.wh)))
}

func (camera *Camera) cornerTransform(x, y float32) (float32, float32) {
	x, y = x-camera.x, y-camera.y
	x, y = -camera.cos*x+camera.sin*y, -camera.sin*x-camera.cos*y
	return camera.x - (x/camera.scale + camera.l), camera.y - (y/camera.scale + camera.t)
}

func (camera *Camera) adjustPosition() {
	w, h := camera.getVisibleArea()
	w2, h2 := w*0.5, h*0.5
	left, right := camera.wl+w2, camera.wl+camera.ww-w2
	top, bottom := camera.wt+h2, camera.wt+camera.wh-h2
	camera.x, camera.y = clamp(camera.x, left, right), clamp(camera.y, top, bottom)
}

func (camera *Camera) adjustScale() {
	rw, rh := camera.getVisibleArea()    // rotated frame: area around the window, rotated without scaling
	sx, sy := rw/camera.ww, rh/camera.wh // vert/horiz scale: minimun scales that the window needs to occupy the world
	rscale := float32(math.Max(float64(sx), float64(sy)))
	camera.scale = float32(math.Max(float64(camera.scale), float64(rscale)))
}

func (camera *Camera) SetWorld(l, t, w, h float32) {
	camera.wl, camera.wt, camera.ww, camera.wh = l, t, w, h
	camera.adjustPosition()
}

func (camera *Camera) SetWindow(l, t, w, h float32) {
	camera.l, camera.t, camera.w, camera.h, camera.w2, camera.h2 = l, t, w, h, w*0.5, h*0.5
	camera.adjustPosition()
}

func (camera *Camera) SetPosition(x, y float32) {
	camera.x, camera.y = x, y
	camera.adjustPosition()
}

func (camera *Camera) SetScale(scale float32) {
	camera.scale = scale
	camera.adjustScale()
	camera.adjustPosition()
}

func (camera *Camera) SetAngle(angle float32) {
	camera.angle = angle
	camera.cos, camera.sin = float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))
	camera.adjustScale()
	camera.adjustPosition()
}

func (camera *Camera) GetWorld() (l, t, w, h float32) {
	return camera.wl, camera.wt, camera.ww, camera.wh
}

func (camera *Camera) GetWindow() (l, t, w, h float32) {
	return camera.l, camera.t, camera.w, camera.h
}

func (camera *Camera) GetPosition() (x, y float32) {
	return camera.x, camera.y
}

func (camera *Camera) GetScale() float32 {
	return camera.scale
}

func (camera *Camera) GetAngle() float32 {
	return camera.angle
}

func (camera *Camera) GetVisible() (l, t, w, h float32) {
	w, h = camera.getVisibleArea()
	return camera.x - w*0.5, camera.y - h*0.5, w, h
}

func (camera *Camera) GetVisibleCorners() (x1, y1, x2, y2, x3, y3, x4, y4 float32) {
	x1, y1 = camera.cornerTransform(camera.x-camera.w2, camera.y-camera.h2)
	x2, y2 = camera.cornerTransform(camera.x+camera.w2, camera.y-camera.h2)
	x3, y3 = camera.cornerTransform(camera.x+camera.w2, camera.y+camera.h2)
	x4, y4 = camera.cornerTransform(camera.x-camera.w2, camera.y+camera.h2)
	return x1, y1, x2, y2, x3, y3, x4, y4
}

func (camera *Camera) Draw(block func(l, t, w, h float32)) {
	l, t, w, h := camera.GetWindow()
	gfx.SetScissor(int32(l), int32(t), int32(w), int32(h))
	gfx.Push()
	gfx.Scale(camera.scale)
	gfx.Translate((camera.w2+camera.l)/camera.scale, (camera.h2+camera.t)/camera.scale)
	gfx.Rotate(-camera.angle)
	gfx.Translate(-camera.x, -camera.y)
	l, t, w, h = camera.GetVisible()
	block(l, t, w, h)
	gfx.Pop()
	gfx.SetScissor()
}

func (camera *Camera) ToWorld(x, y float32) (float32, float32) {
	x, y = (x-camera.w2-camera.l)/camera.scale, (y-camera.h2-camera.t)/camera.scale
	x, y = camera.cos*x-camera.sin*y, camera.sin*x+camera.cos*y
	return x + camera.x, y + camera.y
}

func (camera *Camera) ToScreen(x, y float32) (float32, float32) {
	x, y = x-camera.x, y-camera.y
	x, y = camera.cos*x+camera.sin*y, -camera.sin*x+camera.cos*y
	return camera.scale*x + camera.w2 + camera.l, camera.scale*y + camera.h2 + camera.t
}

func (camera *Camera) Shake(intensity float32) {
	camera.shakeIntensity = float32(math.Min(maxShake, float64(camera.shakeIntensity+intensity)))
}

func (camera *Camera) Update(dt float32) {
	camera.shakeIntensity = float32(math.Max(0, float64(camera.shakeIntensity-atenuationSpeed*dt)))
	if camera.shakeIntensity > 0 {
		x, y := camera.GetPosition()
		x = x + (100-200*randMax(camera.shakeIntensity))*dt
		y = y + (100-200*randMax(camera.shakeIntensity))*dt
		camera.SetPosition(x, y)
	}
}

func clamp(x, minX, maxX float32) float32 {
	if x < minX {
		return minX
	} else if x > maxX {
		return maxX
	}
	return x
}

func randMax(max float32) float32 {
	return randRange(0, max)
}

func randRange(min, max float32) float32 {
	return (rand.Float32() * (max - min)) + min
}
