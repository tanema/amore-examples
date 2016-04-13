package phys

import (
	"fmt"
	"math"

	"github.com/tanema/amore/gfx"
)

type World struct {
	width       float32
	height      float32
	cell_width  float32
	cell_height float32
	numWidth    float32
	numHeight   float32
	grid        [][]*Cell
}

func NewWorld(width, height, cell_size float32) *World {
	cell_width := nextEqualyDividable(width, cell_size)
	cell_height := nextEqualyDividable(height, cell_size)
	numWidth := width / cell_width
	numHeight := height / cell_height

	//initialize grid
	grid := make([][]*Cell, int(numWidth))
	for i := 0; i < int(numWidth); i++ {
		grid[i] = make([]*Cell, int(numHeight))
		for j := 0; j < int(numHeight); j++ {
			grid[i][j] = &Cell{
				x:      float32(i) * cell_width,
				y:      float32(j) * cell_height,
				width:  cell_width,
				height: cell_height,
			}
		}
	}

	return &World{
		width:       width,
		height:      height,
		cell_width:  cell_width,
		cell_height: cell_height,
		numWidth:    numWidth,
		numHeight:   numHeight,
		grid:        grid,
	}
}

func nextEqualyDividable(unit, size float32) float32 {
	for ; math.Mod(float64(unit), float64(size)) != 0; size++ {
	}
	return size
}

func (world *World) AddBody(collidable Collidable, name string, x, y, scale float32, points []float32) *Body {
	new_body := newBody(world, collidable, name, points)
	new_body.Move(x, y, 0, scale)
	return new_body
}

func (world *World) CellAt(x, y float32) *Cell {
	x = float32(math.Floor(float64(x / world.cell_width)))
	y = float32(math.Floor(float64(y / world.cell_height)))
	if x >= world.numWidth {
		x -= world.numWidth
	}
	if x < 0 {
		x = world.numWidth - 1
	}
	if y >= world.numHeight {
		y = 0
	}
	if y < 0 {
		y = world.numHeight - 1
	}
	return world.grid[int(x)][int(y)]
}

func (world *World) DrawGrid() {
	// This really sucks and isn't the best but it is for debugging so it shouldnt
	// be used all the time
	count := map[*Body]bool{}
	for _, row := range world.grid {
		for _, cell := range row {
			for _, body := range cell.bodies {
				count[body] = true
			}
		}
	}

	gfx.Print(fmt.Sprintf("physical objects: %v", len(count)), 0, 30)
	gfx.SetColor(100, 100, 100, 100)
	for i := world.cell_width; i < world.width; i += world.cell_width {
		gfx.Line(i, 0, i, world.height)
	}
	for i := world.cell_height; i < world.height; i += world.cell_height {
		gfx.Line(0, i, world.width, i)
	}
	gfx.SetColor(255, 255, 255, 255)
}
