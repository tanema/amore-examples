package ui

import (
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/mouse"
)

// Button is a clickable rectangle
type Button struct {
	x, y, w, h     float32
	image          *gfx.Image
	imageW, imageH float32
	color          *gfx.Color
	mouseDown      bool
	clickDown      bool
	onClick        func(*Button)
}

// NewButton builds and returns a new button
func NewButton(x, y, w, h float32, image *gfx.Image, color *gfx.Color, onClick func(*Button)) *Button {
	button := &Button{
		x: x, y: y, w: w, h: h,
		color:   color,
		onClick: onClick,
	}
	button.SetImage(image)
	return button
}

// SetImage sets the image for the button
func (button *Button) SetImage(image *gfx.Image) {
	button.image = image
	button.imageW = button.w / float32(image.GetWidth())
	button.imageH = button.h / float32(image.GetHeight())
}

// Update checks for clicks
func (button *Button) Update(dt float32) {
	x, y := mouse.GetPosition()
	isDown := mouse.IsDown(mouse.LeftButton)
	if isDown && !button.mouseDown && button.containsPoint(x, y) {
		button.clickDown = true
	} else if !isDown {
		if button.mouseDown && button.clickDown && button.containsPoint(x, y) {
			button.onClick(button)
		}
		button.clickDown = false
	}
	button.mouseDown = isDown
}

// Draw draws the button
func (button *Button) Draw() {
	gfx.SetColorC(button.color)
	gfx.Rect(gfx.LINE, button.x, button.y, button.w, button.h)
	gfx.SetColor(255, 255, 255, 255)
	gfx.Draw(button.image, button.x, button.y, 0, button.imageW, button.imageH)
}

func (button *Button) containsPoint(px, py float32) bool {
	return button.x < px &&
		button.x+button.w > px &&
		button.y < py &&
		button.y+button.h > py
}
