package game

import (
	"math"

	"github.com/tanema/amore/gfx"
)

type Segment struct {
	index   int
	p1      PointGroup
	p2      PointGroup
	curve   float32
	color   string
	sprites []*Sprite
	cars    []*Car
	looped  bool
}

func newSegment(n int, lastY, newY, curve float32, color string) *Segment {
	return &Segment{
		index: n,
		p1: PointGroup{
			world: Point{
				y: lastY,
				z: float32(n) * segmentLength,
			},
			camera: Point{},
			screen: Point{},
		},
		p2: PointGroup{
			world: Point{
				y: newY,
				z: float32(n+1) * segmentLength,
			},
			camera: Point{},
			screen: Point{},
		},
		curve:   curve,
		sprites: []*Sprite{},
		cars:    []*Car{},
		color:   color,
	}
}

func (segment *Segment) draw(cameraDepth, width, lanes float32) {
	if segment.p1.camera.z <= cameraDepth { // clip by (already rendered) hill
		return
	}

	r1 := segment.rumbleWidth(segment.p1.screen.w, lanes)
	r2 := segment.rumbleWidth(segment.p2.screen.w, lanes)
	l1 := segment.laneMarkerWidth(segment.p1.screen.w, lanes)
	l2 := segment.laneMarkerWidth(segment.p2.screen.w, lanes)

	grasscolor := colors["grass"]
	if segment.color == "odd" {
		grasscolor = grasscolor.Darken(15)
	}
	gfx.SetColorC(grasscolor)
	gfx.Rect(gfx.FILL, 0, segment.p2.screen.y, width, segment.p1.screen.y-segment.p2.screen.y)

	rumblecolor := colors["rumble"]
	if segment.color == "start" || segment.color == "finish" {
		rumblecolor = colors[segment.color]
	} else if segment.color == "odd" {
		rumblecolor = rumblecolor.Darken(15)
	}
	gfx.SetColorC(rumblecolor)
	gfx.Polygon(gfx.FILL, []float32{segment.p1.screen.x - segment.p1.screen.w - r1, segment.p1.screen.y, segment.p1.screen.x - segment.p1.screen.w, segment.p1.screen.y, segment.p2.screen.x - segment.p2.screen.w, segment.p2.screen.y, segment.p2.screen.x - segment.p2.screen.w - r2, segment.p2.screen.y})
	gfx.Polygon(gfx.FILL, []float32{segment.p1.screen.x + segment.p1.screen.w + r1, segment.p1.screen.y, segment.p1.screen.x + segment.p1.screen.w, segment.p1.screen.y, segment.p2.screen.x + segment.p2.screen.w, segment.p2.screen.y, segment.p2.screen.x + segment.p2.screen.w + r2, segment.p2.screen.y})

	roadcolor := colors["road"]
	if segment.color == "start" || segment.color == "finish" {
		roadcolor = colors[segment.color]
	} else if segment.color == "odd" {
		roadcolor = roadcolor.Darken(15)
	}
	gfx.SetColorC(roadcolor)
	gfx.Polygon(gfx.FILL, []float32{segment.p1.screen.x - segment.p1.screen.w, segment.p1.screen.y, segment.p1.screen.x + segment.p1.screen.w, segment.p1.screen.y, segment.p2.screen.x + segment.p2.screen.w, segment.p2.screen.y, segment.p2.screen.x - segment.p2.screen.w, segment.p2.screen.y})

	if segment.color == "even" {
		lanew1 := segment.p1.screen.w * 2 / lanes
		lanew2 := segment.p2.screen.w * 2 / lanes
		lanex1 := segment.p1.screen.x - segment.p1.screen.w + lanew1
		lanex2 := segment.p2.screen.x - segment.p2.screen.w + lanew2
		gfx.SetColorC(colors["lane"])
		for lane := 1; lane < int(lanes); lane++ {
			gfx.Polygon(gfx.FILL, []float32{lanex1 - l1/2, segment.p1.screen.y, lanex1 + l1/2, segment.p1.screen.y, lanex2 + l2/2, segment.p2.screen.y, lanex2 - l2/2, segment.p2.screen.y})
			lanex1 += lanew1
			lanex2 += lanew2
		}
	}
}

func (segment *Segment) rumbleWidth(projectedRoadWidth, lanes float32) float32 {
	return projectedRoadWidth / float32(math.Max(6, 2*float64(lanes)))
}

func (segment *Segment) laneMarkerWidth(projectedRoadWidth, lanes float32) float32 {
	return projectedRoadWidth / float32(math.Max(32, 8*float64(lanes)))
}
