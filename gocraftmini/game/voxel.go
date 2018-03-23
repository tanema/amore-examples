package game

import (
	"github.com/tanema/amore/gfx"
)

// Voxel is the main drawn box
type Voxel struct {
	x, y, z       float32
	width, height float32
	h, s, l       float32
	shine         float32
	relative      bool
}

func newVoxel(x, y, z, width, height, h, s, l float32, relative bool) *Voxel {
	return &Voxel{
		x: x, y: y, z: z,
		width: width, height: height,
		h: h, s: s, l: l,
		relative: relative,
		shine:    1,
	}
}

func (voxel *Voxel) update(world *World, x, y float32) {
	// starting luminance based on distance from player center
	voxel.shine = baseShine + exp(-(x*x+y*y)/playerShineRange*2)
	if world.sin > 0 { // if, daytime add sunlight
		voxel.shine += world.sin * (2 - world.sin) * (1 - voxel.shine)
	}
}

func (voxel *Voxel) draw(camera *Camera, px, py float32) {
	cellSize := float32(1)
	if voxel.relative {
		cellSize = camera.getCellSize()
	}

	gfx.SetColorC(gfx.NewHSLColor(voxel.h, voxel.s, pow(voxel.l, 1/voxel.shine), 1))
	x, y := camera.worldToScreen(px, py, voxel.z, voxel.relative)
	width, height := voxel.width*cellSize, voxel.height*cellSize
	coords := []float32{
		x + width, y, x, y + width/2,
		x - width, y,
		x - width, y - height,
		x, y - height - width/2,
		x + width, y - height,
	}
	gfx.Polygon(gfx.FILL, coords)
}
