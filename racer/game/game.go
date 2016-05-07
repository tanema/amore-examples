package game

import (
	"math"
	"math/rand"

	"github.com/tanema/amore/gfx"
)

const (
	shortLength          = 0
	mediumLength         = 50
	longLength           = 100
	lowHill              = 20
	mediumHill           = 40
	highHill             = 60
	easyCurve    float32 = 2
	mediumCurve  float32 = 4
	hardCurve    float32 = 6
)

type (
	Point struct {
		x, y, z, w, scale float32
	}
	PointGroup struct {
		world  Point
		camera Point
		screen Point
	}
)

var (
	width         float32 // logical canvas width
	height        float32 // logical canvas height
	player        *Player
	centrifugal   float32    = 0.3                                                       // centrifugal force multiplier when going around curves
	skySpeed      float32    = 0.001                                                     // background sky layer scroll speed when going around curve (or up hill)
	hillSpeed     float32    = 0.002                                                     // background hill layer scroll speed when going around curve (or up hill)
	treeSpeed     float32    = 0.003                                                     // background tree layer scroll speed when going around curve (or up hill)
	skyOffset     float32    = 0                                                         // current sky scroll offset
	hillOffset    float32    = 0                                                         // current hill scroll offset
	treeOffset    float32    = 0                                                         // current tree scroll offset
	segments      []*Segment                                                             // array of road segments
	cars          []*Car                                                                 // array of cars on the road
	resolution    float32    = height / 480                                              // scaling factor to provide resolution independence (computed)
	roadWidth     float32    = 2000                                                      // actually half the roads width, easier math if the road spans from -roadWidth to +roadWidth
	segmentLength float32    = 200                                                       // length of a single segment
	rumbleLength  int        = 3                                                         // number of segments per red/white rumble strip
	trackLength   float32                                                                // z length of entire track (computed)
	lanes         float32    = 3                                                         // number of lanes
	fieldOfView   float32    = 100                                                       // angle (degrees) for field of view
	cameraHeight  float32    = 1000                                                      // z height of camera
	cameraDepth   float32    = float32(1 / math.Tan(float64(fieldOfView/2)*math.Pi/180)) // z distance camera is from screen (computed)
	drawDistance  float32    = 300                                                       // number of segments to draw
	fogDensity    float32    = 8                                                         // exponential fog density
	maxSpeed                 = segmentLength / (1.0 / 60.0)                              // top speed (ensure we can't move more than 1 segment in a single frame to make collision detection easier)
	accel                    = maxSpeed / 5                                              // acceleration rate - tuned until it 'felt' right
	breaking                 = -maxSpeed                                                 // deceleration rate when braking
	decel                    = -maxSpeed / 5                                             // 'natural' deceleration rate when neither accelerating, nor braking
	offRoadDecel             = -maxSpeed / 2                                             // off road deceleration is somewhere in between
	offRoadLimit             = maxSpeed / 4                                              // limit when off road deceleration no longer applies (e.g. you can always go at least this speed even when off road)
	totalCars                = 200                                                       // total number of cars on the road
)

func New() {
	width = gfx.GetWidth()
	height = gfx.GetHeight()
	initSpriteSheets()
	resetRoad()
	player = newPlayer(cameraHeight * cameraDepth)
}

func Update(dt float32) {
	startPosition := player.position
	player.update(dt)
	for _, car := range cars {
		car.update(dt, player)
	}
	skyOffset = increase(skyOffset, skySpeed*player.segment.curve*(player.position-startPosition)/segmentLength, 1)
	hillOffset = increase(hillOffset, hillSpeed*player.segment.curve*(player.position-startPosition)/segmentLength, 1)
	treeOffset = increase(treeOffset, treeSpeed*player.segment.curve*(player.position-startPosition)/segmentLength, 1)
}

func Draw() {
	baseSegment := findSegment(player.position)
	basePercent := percentRemaining(player.position, segmentLength)
	player.segment = findSegment(player.position + player.point.z)
	playerPercent := percentRemaining(player.position+player.point.z, segmentLength)
	playerY := interpolate(player.segment.p1.world.y, player.segment.p2.world.y, playerPercent)

	// render curves from from to back so that we can just render from back to from and
	// not calculate clipping
	var x float32 = 0
	var dx = -(baseSegment.curve * basePercent)
	curves := [][]float32{}
	for n := 0; n < int(drawDistance); n++ {
		segment := segments[(baseSegment.index+n)%len(segments)]
		curves = append(curves, []float32{(player.point.x * roadWidth) - x, (player.point.x * roadWidth) - x - dx})
		x = x + dx
		dx = dx + segment.curve
	}

	drawBackground(width, height, backgrounds["sky"], skyOffset, resolution*skySpeed*playerY)
	drawBackground(width, height, backgrounds["hills"], hillOffset, resolution*hillSpeed*playerY)
	drawBackground(width, height, backgrounds["trees"], treeOffset, resolution*treeSpeed*playerY)

	for n := int(drawDistance - 1); n > 0; n-- {
		segment := segments[(baseSegment.index+n)%len(segments)]
		segment.looped = segment.index < baseSegment.index

		loop := trackLength
		if !segment.looped {
			loop = 0
		}

		project(&segment.p1, curves[n][0], playerY+cameraHeight, player.position-loop, cameraDepth, width, height, roadWidth)
		project(&segment.p2, curves[n][1], playerY+cameraHeight, player.position-loop, cameraDepth, width, height, roadWidth)

		if segment.p1.camera.z > cameraDepth { // clip by (already rendered) hill
			segment.draw(cameraDepth, width, lanes)

			for i := 0; i < len(segment.cars); i++ {
				car := segment.cars[i]
				spriteScale := interpolate(segment.p1.screen.scale, segment.p2.screen.scale, car.percent)
				spriteX := interpolate(segment.p1.screen.x, segment.p2.screen.x, car.percent) + (spriteScale * car.offset * roadWidth * width / 2)
				spriteY := interpolate(segment.p1.screen.y, segment.p2.screen.y, car.percent)
				car.draw(width, height, roadWidth, spriteScale, spriteX, spriteY, -0.5, -1)
			}

			for i := 0; i < len(segment.sprites); i++ {
				sprite := segment.sprites[i]
				spriteScale := segment.p1.screen.scale
				spriteX := segment.p1.screen.x + (spriteScale * sprite.offset * roadWidth * width / 2)
				spriteY := segment.p1.screen.y
				ov := float32(0)
				if sprite.offset < 0 {
					ov = -1
				}
				sprite.draw(width, height, roadWidth, spriteScale, spriteX, spriteY, ov, -1)
			}

			drawFog(0, segment.p1.screen.y, width, segment.p2.screen.y-segment.p1.screen.y, exponentialFog(float32(n)/drawDistance, fogDensity))
		}

	}

	player.draw()
}

func findSegment(z float32) *Segment {
	return segments[int(math.Floor(float64(z/segmentLength)))%len(segments)]
}

func lastY() float32 {
	if len(segments) == 0 {
		return 0
	}
	return segments[len(segments)-1].p2.world.y
}

func addSegment(curve, y float32) {
	var n = len(segments)
	color := "odd"
	if int(math.Floor(float64(float32(n/rumbleLength))))%2 == 0 {
		color = "even"
	}
	segments = append(segments, newSegment(n, lastY(), y, curve, color))
}

func addSprite(n int, sprite *gfx.Quad, offset float32) {
	segments[n].sprites = append(segments[n].sprites, newSprite(sprite, offset))
}

func addRoad(enter, hold, leave int, curve, y float32) {
	startY := lastY()
	endY := startY + (y * segmentLength)
	total := enter + hold + leave
	for n := 0; n < enter; n++ {
		addSegment(easeIn(0, curve, float32(n)/float32(enter)), easeInOut(startY, endY, float32(n)/float32(total)))
	}
	for n := 0; n < hold; n++ {
		addSegment(curve, easeInOut(startY, endY, float32(enter+n)/float32(total)))
	}
	for n := 0; n < leave; n++ {
		addSegment(easeInOut(curve, 0, float32(n)/float32(leave)), easeInOut(startY, endY, float32(enter+hold+n)/float32(total)))
	}
}

func addStraight(num int)                     { addRoad(num, num, num, 0, 0) }
func addHill(num int, height float32)         { addRoad(num, num, num, 0, height) }
func addCurve(num int, curve, height float32) { addRoad(num, num, num, curve, height) }

func addLowRollingHills(num int, height float32) {
	addRoad(num, num, num, 0, height/2)
	addRoad(num, num, num, 0, -height)
	addRoad(num, num, num, easyCurve, height)
	addRoad(num, num, num, 0, 0)
	addRoad(num, num, num, -easyCurve, height/2)
	addRoad(num, num, num, 0, 0)
}

func addSCurves() {
	addRoad(mediumLength, mediumLength, mediumLength, -easyCurve, 0)
	addRoad(mediumLength, mediumLength, mediumLength, mediumCurve, mediumHill)
	addRoad(mediumLength, mediumLength, mediumLength, easyCurve, -lowHill)
	addRoad(mediumLength, mediumLength, mediumLength, -easyCurve, mediumHill)
	addRoad(mediumLength, mediumLength, mediumLength, -mediumCurve, -mediumHill)
}

func addBumps() {
	addRoad(10, 10, 10, 0, 5)
	addRoad(10, 10, 10, 0, -2)
	addRoad(10, 10, 10, 0, -5)
	addRoad(10, 10, 10, 0, 8)
	addRoad(10, 10, 10, 0, 5)
	addRoad(10, 10, 10, 0, -7)
	addRoad(10, 10, 10, 0, 5)
	addRoad(10, 10, 10, 0, -2)
}

func addDownhillToEnd(num int) {
	addRoad(num, num, num, -easyCurve, -lastY()/segmentLength)
}

func resetRoad() {
	segments = []*Segment{}

	addStraight(shortLength)
	addLowRollingHills(shortLength, lowHill)
	addSCurves()
	addCurve(mediumLength, mediumCurve, lowHill)
	addBumps()
	addLowRollingHills(shortLength, lowHill)
	addCurve(longLength*2, mediumCurve, mediumHill)
	addStraight(mediumLength)
	addHill(mediumLength, highHill)
	addSCurves()
	addCurve(longLength, -mediumCurve, 0)
	addHill(longLength, highHill)
	addCurve(longLength, mediumCurve, -lowHill)
	addBumps()
	addHill(longLength, -mediumHill)
	addStraight(mediumLength)
	addSCurves()
	addDownhillToEnd(200)

	resetSprites()
	resetCars()

	segments[10].color = "start"
	segments[11].color = "start"
	for n := 0; n < rumbleLength; n++ {
		segments[len(segments)-1-n].color = "finish"
	}

	trackLength = float32(len(segments) * int(segmentLength))
}

func resetSprites() {
	addSprite(20, spriteSheet["billboard07"], -1)
	addSprite(40, spriteSheet["billboard06"], -1)
	addSprite(60, spriteSheet["billboard08"], -1)
	addSprite(80, spriteSheet["billboard09"], -1)
	addSprite(100, spriteSheet["billboard01"], -1)
	addSprite(120, spriteSheet["billboard02"], -1)
	addSprite(140, spriteSheet["billboard03"], -1)
	addSprite(160, spriteSheet["billboard04"], -1)
	addSprite(180, spriteSheet["billboard05"], -1)
	addSprite(240, spriteSheet["billboard07"], -1.2)
	addSprite(240, spriteSheet["billboard06"], 1.2)
	addSprite(len(segments)-25, spriteSheet["billboard07"], -1.2)
	addSprite(len(segments)-25, spriteSheet["billboard06"], 1.2)

	for n := 10; n < 200; n += 4 + int(math.Floor(float64(n)/100)) {
		addSprite(n, spriteSheet["palm_tree"], randRange(0.5, 1))
		addSprite(n, spriteSheet["palm_tree"], randRange(1, 3))
	}

	for n := 250; n < 1000; n += 5 {
		addSprite(n, spriteSheet["column"], 1.1)
		addSprite(n+rand.Intn(5), spriteSheet["tree1"], -1-randMax(2))
		addSprite(n+rand.Intn(5), spriteSheet["tree2"], -1-randMax(2))
	}

	for n := 200; n < len(segments); n += 3 {
		choice := rand.Intn(len(plants) - 1)
		side := -1
		if rand.Intn(1) == 0 {
			side = 1
		}
		addSprite(n, plants[choice], float32(side)*randRange(2, 7))
	}

	for n := 1000; n < (len(segments) - 50); n += 100 {
		side := float32(-1)
		if rand.Intn(1) == 0 {
			side = 1
		}
		choice := rand.Intn(len(billboards) - 1)
		addSprite(n+rand.Intn(50), billboards[choice], side)
		for i := 0; i < 20; i++ {
			side = -1
			if rand.Intn(1) == 0 {
				side = 1
			}
			choice = rand.Intn(len(plants) - 1)
			addSprite(n+rand.Intn(50), plants[choice], side*(1.5+rand.Float32()))
		}
	}
}

func resetCars() {
	cars = []*Car{}
	for n := 0; n < totalCars; n++ {
		offset := rand.Float32()
		if rand.Intn(1) == 0 {
			offset *= -0.8
		} else {
			offset *= 0.8
		}

		z := float32(math.Floor(rand.Float64()*float64(len(segments)))) * segmentLength
		choice := rand.Intn(len(carQuads) - 1)
		sprite := carQuads[choice]
		speed := maxSpeed/4 + rand.Float32()*maxSpeed/2
		if sprite == spriteSheet["semi"] {
			speed = maxSpeed/4 + rand.Float32()*maxSpeed/4
		}

		segment := findSegment(z)
		car := newCar(segment, sprite, z, offset, speed)
		segment.cars = append(segment.cars, car)
		cars = append(cars, car)
	}
}

func project(p *PointGroup, cameraX, cameraY, cameraZ, cameraDepth, width, height, roadWidth float32) {
	p.camera.x = p.world.x - cameraX
	p.camera.y = p.world.y - cameraY
	p.camera.z = p.world.z - cameraZ
	p.screen.scale = cameraDepth / p.camera.z
	p.screen.x = round((width / 2) + (p.screen.scale * p.camera.x * width / 2))
	p.screen.y = round((height / 2) - (p.screen.scale * p.camera.y * height / 2))
	p.screen.w = round((p.screen.scale * roadWidth * width / 2))
}

func drawBackground(width, height float32, layer *gfx.Quad, rotation, offset float32) {
	layerx, layery, layerw, layerh := layer.GetViewport()
	imagew := layerw / 2
	imageh := layerh

	sourcex := layerx + int32(math.Floor(float64(float32(layerw)*rotation)))
	sourcey := layery
	sourcew := int32(math.Min(float64(imagew), float64(layerx+layerw-sourcex)))
	sourceh := imageh

	desty := offset
	destw := float32(math.Floor(float64(width * float32(sourcew/imagew))))
	desth := height

	sourceq := gfx.NewQuad(sourcex, sourcey, sourcew, sourceh, background.GetWidth(), background.GetHeight())
	gfx.Drawq(background, sourceq, 0, desty, 0, float32(layerw)/destw, float32(layerh)/desth)
	if sourcew < imagew {
		sourceq = gfx.NewQuad(layerx, sourcey, imagew, sourceh, background.GetWidth(), background.GetHeight())
		gfx.Drawq(background, sourceq, destw-1, desty, 0, float32(layerw)/(width-destw), float32(layerh)/desth)
	}
}

func drawFog(x, y, width, height, fog float32) {
	if fog < 1 {
		gfx.SetColor(0, 81, 8, (1-fog)*255)
		gfx.Rect(gfx.FILL, x, y, width, height)
	}
}
