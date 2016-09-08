package game

import (
	"github.com/tanema/amore/gfx"
)

const (
	maxShake        = float32(5)
	atenuationSpeed = float32(4)
)

type Camera struct {
	x, y           float32
	width, height  float32
	halfWidth      float32
	halfHeight     float32
	scale          float32
	rot            float32
	shakeIntensity float32
}

func NewCamera() *Camera {
	w, h := gfx.GetWidth(), gfx.GetHeight()

	return &Camera{
		width:      w,
		height:     h,
		halfWidth:  w * 0.5,
		halfHeight: h * 0.5,
		scale:      1,
	}
}

func (camera *Camera) LookAt(x, y float32) {
	camera.x, camera.y = x, y
}

func (camera *Camera) Move(dx, dy float32) {
	camera.x, camera.y = camera.x+dx, camera.y+dy
}

func (camera *Camera) Rotate(phi float32) {
	camera.rot = camera.rot + phi
}

func (camera *Camera) RotateTo(phi float32) {
	camera.rot = phi
}

func (camera *Camera) Zoom(mul float32) {
	camera.scale = camera.scale * mul
}

func (camera *Camera) ZoomTo(zoom float32) {
	camera.scale = zoom
}

func (camera *Camera) GetVisible() (l, t, w, h float32) {
	return camera.x - camera.halfWidth, camera.y - camera.halfHeight, camera.width, camera.height
}

func (camera *Camera) Draw(draw func(l, t, w, h float32)) {
	gfx.Push()
	{
		gfx.Translate(-camera.x, -camera.y)
		gfx.Scale(camera.scale)
		gfx.Rotate(camera.rot)
		gfx.Translate(camera.halfWidth, camera.halfHeight)

		l, t, w, h := camera.GetVisible()
		draw(l, t, w, h)
	}
	gfx.Pop()
}

func (camera *Camera) Shake(intensity float32) {
	camera.shakeIntensity = min(maxShake, camera.shakeIntensity+intensity)
}

func (camera *Camera) Update(dt float32) {
	camera.shakeIntensity = max(0, camera.shakeIntensity-atenuationSpeed*dt)
	if camera.shakeIntensity > 0 {
		camera.x += (100 - 200*randMax(camera.shakeIntensity)) * dt
		camera.y += (100 - 200*randMax(camera.shakeIntensity)) * dt
	}
}
