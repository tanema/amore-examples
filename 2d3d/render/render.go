package render

import (
	"github.com/tanema/amore/gfx"
)

const increment float32 = 1

type Fake3D struct {
	img    *gfx.Image
	quads  []*gfx.Quad
	ox, oy float32
}

func New(filepath string, frameWidth, frameHeight int32) (*Fake3D, error) {
	img, err := gfx.NewImage(filepath)
	if err != nil {
		return nil, err
	}

	imageWidth, imageHeight := img.GetWidth(), img.GetHeight()
	framesWide, framesHeight := imageWidth/frameWidth, imageHeight/frameHeight
	quads := make([]*gfx.Quad, 0, framesWide*framesHeight)

	var x, y int32
	for y = 0; y < framesHeight; y++ {
		for x = 0; x < framesWide; x++ {
			newQuad := gfx.NewQuad(x*frameWidth, y*frameHeight, frameWidth, frameHeight, imageWidth, imageHeight)
			quads = append(quads, newQuad)
		}
	}

	return &Fake3D{
		img:   img,
		quads: quads,
		ox:    float32(frameWidth) / 2,
		oy:    float32(frameHeight) / 2,
	}, nil
}

func (f3d *Fake3D) Draw(x, y, angle float32) {
	for _, quad := range f3d.quads {
		gfx.Drawq(f3d.img, quad, x, y, angle, 1, 1, f3d.ox, f3d.oy)
		y -= increment
	}
}
