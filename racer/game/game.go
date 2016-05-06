package game

import (
	"math"
	"math/rand"

	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"
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
	Segment struct {
		index   int
		p1      PointGroup
		p2      PointGroup
		curve   float32
		color   string
		sprites []*Sprite
		cars    []*Car
		looped  bool
		fog     float32
		clip    float32
	}
	Car struct {
		offset  float32
		percent float32
		point   Point
		sprite  *Sprite
		speed   float32
	}
	Sprite struct {
		source *gfx.Quad
		offset float32
	}
)

var (
	width         float32                                                                // logical canvas width
	height        float32                                                                // logical canvas height
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
	playerX       float32    = 0                                                         // player x offset from center of road (-1 to 1 to stay independent of roadWidth)
	playerZ       float32                                                                // player relative z distance from camera (computed)
	fogDensity    float32    = 5                                                         // exponential fog density
	position      float32    = 0                                                         // current camera Z position (add playerZ to get player's absolute Z position)
	speed         float32    = 0                                                         // current speed
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
}

func Update(dt float32) {
	var playerSegment = findSegment(position + playerZ)
	var playerW = spriteSheet["player_straight"].GetWidth() * sprite_scale
	var speedPercent = speed / maxSpeed
	var dx = dt * 2 * speedPercent // at top speed, should be able to cross from left to right (-1 to 1) in 1 second
	var startPosition = position

	updateCars(dt, playerSegment, playerW)

	position = increase(position, dt*speed, trackLength)

	if keyboard.IsDown(keyboard.KeyLeft) {
		playerX = playerX - dx
	} else if keyboard.IsDown(keyboard.KeyRight) {
		playerX = playerX + dx
	}

	playerX = playerX - (dx * speedPercent * playerSegment.curve * centrifugal)

	if keyboard.IsDown(keyboard.KeyUp) {
		speed = accelerate(speed, accel, dt)
	} else if keyboard.IsDown(keyboard.KeyDown) {
		speed = accelerate(speed, breaking, dt)
	} else {
		speed = accelerate(speed, decel, dt)
	}

	if (playerX < -1) || (playerX > 1) {

		if speed > offRoadLimit {
			speed = accelerate(speed, offRoadDecel, dt)
		}

		for n := 0; n < len(playerSegment.sprites); n++ {
			sprite := playerSegment.sprites[n]
			spriteW := sprite.source.GetWidth() * sprite_scale

			ov := float32(-1)
			if sprite.offset > 0 {
				ov = 1
			}
			if overlap(playerX, playerW, sprite.offset+spriteW/2*ov, spriteW, 1) {
				speed = maxSpeed / 5
				position = increase(playerSegment.p1.world.z, -playerZ, trackLength) // stop in front of sprite (at front of segment)
				break
			}
		}
	}

	for n := 0; n < len(playerSegment.cars); n++ {
		car := playerSegment.cars[n]
		carW := car.sprite.source.GetWidth() * sprite_scale
		if speed > car.speed {
			if overlap(playerX, playerW, car.offset, carW, 0.8) {
				speed = car.speed * (car.speed / speed)
				position = increase(car.point.z, -playerZ, trackLength)
				break
			}
		}
	}

	playerX = clamp(playerX, -3, 3)   // dont ever let it go too far out of bounds
	speed = clamp(speed, 0, maxSpeed) // or exceed maxSpeed

	skyOffset = increase(skyOffset, skySpeed*playerSegment.curve*(position-startPosition)/segmentLength, 1)
	hillOffset = increase(hillOffset, hillSpeed*playerSegment.curve*(position-startPosition)/segmentLength, 1)
	treeOffset = increase(treeOffset, treeSpeed*playerSegment.curve*(position-startPosition)/segmentLength, 1)
}

func updateCars(dt float32, playerSegment *Segment, playerW float32) {
	for n := 0; n < len(cars); n++ {
		car := cars[n]
		oldSegment := findSegment(car.point.z)
		car.offset = car.offset + updateCarOffset(car, oldSegment, playerSegment, playerW)
		car.point.z = increase(car.point.z, dt*car.speed, trackLength)
		car.percent = percentRemaining(car.point.z, segmentLength) // useful for interpolation during rendering phase
		newSegment := findSegment(car.point.z)
		if oldSegment != newSegment {
			for i, c := range oldSegment.cars {
				if car == c {
					oldSegment.cars = append(oldSegment.cars[:i], oldSegment.cars[i+1:]...)
					break
				}
			}
			newSegment.cars = append(newSegment.cars, car)
		}
	}
}

func updateCarOffset(car *Car, carSegment, playerSegment *Segment, playerW float32) float32 {
	var dir float32
	lookahead := 20
	carW := car.sprite.source.GetWidth() * sprite_scale

	// optimization, dont bother steering around other cars when 'out of sight' of the player
	if float32(carSegment.index-playerSegment.index) > drawDistance {
		return 0
	}

	for i := 1; i < lookahead; i++ {
		segment := segments[(carSegment.index+i)%len(segments)]

		if (segment == playerSegment) && (car.speed > speed) && (overlap(playerX, playerW, car.offset, carW, 1.2)) {
			if playerX > 0.5 {
				dir = -1
			} else if (playerX < -0.5) || car.offset > playerX {
				dir = 1
			} else {
				dir = -1
			}
			return dir * 1 / float32(i) * (car.speed - speed) / maxSpeed // the closer the cars (smaller i) and the greated the speed ratio, the larger the offset
		}

		for j := 0; j < len(segment.cars); j++ {
			otherCar := segment.cars[j]
			otherCarW := otherCar.sprite.source.GetWidth() * sprite_scale
			if (car.speed > otherCar.speed) && overlap(car.offset, carW, otherCar.offset, otherCarW, 1.2) {
				if otherCar.offset > 0.5 {
					dir = -1
				} else if (otherCar.offset < -0.5) || (car.offset > otherCar.offset) {
					dir = 1
				} else {
					dir = -1
				}
				return dir * 1 / float32(i) * (car.speed - otherCar.speed) / maxSpeed
			}
		}
	}

	// if no cars ahead, but I have somehow ended up off road, then steer back on
	if car.offset < -0.9 {
		return 0.1
	} else if car.offset > 0.9 {
		return -0.1
	}
	return 0
}

func Draw() {

	var baseSegment = findSegment(position)
	var basePercent = percentRemaining(position, segmentLength)
	var playerSegment = findSegment(position + playerZ)
	var playerPercent = percentRemaining(position+playerZ, segmentLength)
	var playerY = interpolate(playerSegment.p1.world.y, playerSegment.p2.world.y, playerPercent)
	var maxy = height

	var x float32 = 0
	var dx = -(baseSegment.curve * basePercent)

	render_background(width, height, backgrounds["sky"], skyOffset, resolution*skySpeed*playerY)
	render_background(width, height, backgrounds["hills"], hillOffset, resolution*hillSpeed*playerY)
	render_background(width, height, backgrounds["trees"], treeOffset, resolution*treeSpeed*playerY)

	for n := 0; n < int(drawDistance); n++ {
		segment := segments[(baseSegment.index+n)%len(segments)]
		segment.looped = segment.index < baseSegment.index
		segment.fog = exponentialFog(float32(n)/drawDistance, fogDensity)
		segment.clip = maxy

		loop := trackLength
		if !segment.looped {
			loop = 0
		}

		project(&segment.p1, (playerX*roadWidth)-x, playerY+cameraHeight, position-loop, cameraDepth, width, height, roadWidth)
		project(&segment.p2, (playerX*roadWidth)-x-dx, playerY+cameraHeight, position-loop, cameraDepth, width, height, roadWidth)

		x = x + dx
		dx = dx + segment.curve

		if (segment.p1.camera.z <= cameraDepth) || // behind us
			(segment.p2.screen.y >= segment.p1.screen.y) || // back face cull
			(segment.p2.screen.y >= maxy) { // clip by (already rendered) hill
			continue
		}

		render_segment(width, lanes,
			segment.p1.screen.x,
			segment.p1.screen.y,
			segment.p1.screen.w,
			segment.p2.screen.x,
			segment.p2.screen.y,
			segment.p2.screen.w,
			segment.fog,
			segment.color)

		maxy = segment.p1.screen.y
	}

	for n := (drawDistance - 1); n > 0; n-- {
		segment := segments[(baseSegment.index+int(n))%len(segments)]

		for i := 0; i < len(segment.cars); i++ {
			car := segment.cars[i]
			spriteScale := interpolate(segment.p1.screen.scale, segment.p2.screen.scale, car.percent)
			spriteX := interpolate(segment.p1.screen.x, segment.p2.screen.x, car.percent) + (spriteScale * car.offset * roadWidth * width / 2)
			spriteY := interpolate(segment.p1.screen.y, segment.p2.screen.y, car.percent)
			render_sprite(width, height, resolution, roadWidth, car.sprite.source, spriteScale, spriteX, spriteY, -0.5, -1, segment.clip)
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
			render_sprite(width, height, resolution, roadWidth, sprite.source, spriteScale, spriteX, spriteY, ov, -1, segment.clip)
		}

		lr := float32(0)
		if keyboard.IsDown(keyboard.KeyLeft) {
			lr = -1
		} else if keyboard.IsDown(keyboard.KeyRight) {
			lr = 1
		}
		if segment == playerSegment {
			render_player(width, height, resolution, roadWidth, speed/maxSpeed,
				cameraDepth/playerZ,
				width/2,
				(height/2)-(cameraDepth/playerZ*interpolate(playerSegment.p1.camera.y, playerSegment.p2.camera.y, playerPercent)*height/2),
				speed*lr,
				playerSegment.p2.world.y-playerSegment.p1.world.y)
		}
	}
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
	segments = append(segments, &Segment{
		index: n,
		p1: PointGroup{
			world: Point{
				y: lastY(),
				z: float32(n) * segmentLength,
			},
			camera: Point{},
			screen: Point{},
		},
		p2: PointGroup{
			world: Point{
				y: y,
				z: float32(n+1) * segmentLength,
			},
			camera: Point{},
			screen: Point{},
		},
		curve:   curve,
		sprites: []*Sprite{},
		cars:    []*Car{},
		color:   color,
	})
}

func addSprite(n int, sprite *gfx.Quad, offset float32) {
	segments[n].sprites = append(segments[n].sprites, &Sprite{
		source: sprite,
		offset: offset,
	})
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
	playerZ = (cameraHeight * cameraDepth)
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

	segments[findSegment(playerZ).index+2].color = "start"
	segments[findSegment(playerZ).index+3].color = "start"
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

		car := &Car{
			offset: offset,
			point: Point{
				z: z,
			},
			sprite: &Sprite{
				source: sprite,
			},
			speed: speed,
		}

		segment := findSegment(car.point.z)
		segment.cars = append(segment.cars, car)
		cars = append(cars, car)
	}
}

func accelerate(v, accel, dt float32) float32 {
	val := v + (accel * dt)
	return val
}

func easeIn(a, b, percent float32) float32 {
	return a + (b-a)*float32(math.Pow(float64(percent), 2))
}

func easeOut(a, b, percent float32) float32 {
	return a + (b-a)*(1-float32(math.Pow(1-float64(percent), 2)))
}

func easeInOut(a, b, percent float32) float32 {
	return a + (b-a)*((-float32(math.Cos(float64(percent)*math.Pi))/2)+0.5)
}

func interpolate(a, b, percent float32) float32 {
	return a + (b-a)*percent
}

func clamp(value, min, max float32) float32 {
	return float32(math.Max(float64(min), math.Min(float64(value), float64(max))))
}

func exponentialFog(distance, density float32) float32 {
	return 1 / float32(math.Pow(math.E, float64(distance*distance*density)))
}

func percentRemaining(n, total float32) float32 {
	return float32(math.Mod(float64(n), float64(total))) / total
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

func overlap(x1, w1, x2, w2, percent float32) bool {
	var half = percent / 2
	var min1 = x1 - (w1 * half)
	var max1 = x1 + (w1 * half)
	var min2 = x2 - (w2 * half)
	var max2 = x2 + (w2 * half)
	return !((max1 < min2) || (min1 > max2))
}

func increase(start, increment, max float32) float32 { // with looping
	var result = start + increment
	if math.IsInf(float64(increment), 1) {
		panic("ahhh")
	}
	for result >= max {
		result -= max
	}
	for result < 0 {
		result += max
	}
	return result
}

func round(val float32) float32 {
	if val < 0 {
		return float32(int(val - 0.5))
	}
	return float32(int(val + 0.5))
}

func render_segment(width, lanes, x1, y1, w1, x2, y2, w2, fog float32, color string) {
	r1 := rumbleWidth(w1, lanes)
	r2 := rumbleWidth(w2, lanes)
	l1 := laneMarkerWidth(w1, lanes)
	l2 := laneMarkerWidth(w2, lanes)

	grassColor := colors["grass"]
	if color == "odd" {
		grassColor = grassColor.Darken(15)
	}
	gfx.SetColorC(grassColor)
	gfx.Rect(gfx.FILL, 0, y2, width, y1-y2)

	rumbleColor := colors["rumble"]
	if color == "start" || color == "finish" {
		rumbleColor = colors[color]
	} else if color == "odd" {
		rumbleColor = rumbleColor.Darken(15)
	}
	gfx.SetColorC(rumbleColor)
	gfx.Polygon(gfx.FILL, []float32{x1 - w1 - r1, y1, x1 - w1, y1, x2 - w2, y2, x2 - w2 - r2, y2})
	gfx.Polygon(gfx.FILL, []float32{x1 + w1 + r1, y1, x1 + w1, y1, x2 + w2, y2, x2 + w2 + r2, y2})

	roadColor := colors["road"]
	if color == "start" || color == "finish" {
		roadColor = colors[color]
	} else if color == "odd" {
		roadColor = roadColor.Darken(15)
	}
	gfx.SetColorC(roadColor)
	gfx.Polygon(gfx.FILL, []float32{x1 - w1, y1, x1 + w1, y1, x2 + w2, y2, x2 - w2, y2})

	if color == "even" {
		lanew1 := w1 * 2 / lanes
		lanew2 := w2 * 2 / lanes
		lanex1 := x1 - w1 + lanew1
		lanex2 := x2 - w2 + lanew2
		gfx.SetColorC(colors["lane"])
		for lane := 1; lane < int(lanes); lane++ {
			gfx.Polygon(gfx.FILL, []float32{lanex1 - l1/2, y1, lanex1 + l1/2, y1, lanex2 + l2/2, y2, lanex2 - l2/2, y2})
			lanex1 += lanew1
			lanex2 += lanew2
		}
	}

	render_fog(0, y1, width, y2-y1, fog)
}

func render_background(width, height float32, layer *gfx.Quad, rotation, offset float32) {
	layerx, layery, layerw, layerh := layer.GetViewport()
	imageW := layerw / 2
	imageH := layerh

	sourceX := layerx + int32(math.Floor(float64(float32(layerw)*rotation)))
	sourceY := layery
	sourceW := int32(math.Min(float64(imageW), float64(layerx+layerw-sourceX)))
	sourceH := imageH

	destY := offset
	destW := float32(math.Floor(float64(width * float32(sourceW/imageW))))
	destH := height

	sourceq := gfx.NewQuad(sourceX, sourceY, sourceW, sourceH, background.GetWidth(), background.GetHeight())
	gfx.Drawq(background, sourceq, 0, destY, 0, float32(layerw)/destW, float32(layerh)/destH)
	if sourceW < imageW {
		sourceq = gfx.NewQuad(layerx, sourceY, imageW, sourceH, background.GetWidth(), background.GetHeight())
		gfx.Drawq(background, sourceq, destW-1, destY, 0, float32(layerw)/(width-destW), float32(layerh)/destH)
	}
}

func render_sprite(width, height, resolution, roadWidth float32, sprite *gfx.Quad, scale, destX, destY, offsetX, offsetY, clipY float32) {
	//  scale for projection AND relative to roadWidth (for tweakUI)
	spritex, spritey, spritew, spriteh := sprite.GetViewport()
	destW := (float32(spritew) * scale * width / 2) * (sprite_scale * roadWidth)
	destH := (float32(spriteh) * scale * width / 2) * (sprite_scale * roadWidth)

	destX = destX + (destW * offsetX)
	destY = destY + (destH * offsetY)

	clipH := float32(0)
	if clipY != 0 {
		clipH = float32(math.Max(0, float64(destY+destH-clipY)))
	}
	if clipH < destH {
		gfx.SetColor(255, 255, 255, 255)
		sourceq := gfx.NewQuad(spritex, spritey, spritew, spriteh-int32(float32(spriteh)*clipH/destH), sprites.GetWidth(), sprites.GetHeight())
		gfx.Drawq(sprites, sourceq, destX, destY, 0, destW/float32(spritew), (destH-clipH)/float32(spriteh))
	}
}

func render_player(width, height, resolution, roadWidth, speedPercent, scale, destX, destY, steer, updown float32) {
	side := float32(-1)
	if rand.Intn(1) == 0 {
		side = 1
	}
	bounce := (1.5 * rand.Float32() * speedPercent * resolution) * side

	var sprite *gfx.Quad
	if steer < 0 {
		if updown > 0 {
			sprite = spriteSheet["player_uphill_left"]
		} else {
			sprite = spriteSheet["player_left"]
		}
	} else if steer > 0 {
		if updown > 0 {
			sprite = spriteSheet["player_uphill_right"]
		} else {
			sprite = spriteSheet["player_right"]
		}
	} else {
		if updown > 0 {
			sprite = spriteSheet["player_uphill_straight"]
		} else {
			sprite = spriteSheet["player_straight"]
		}
	}

	render_sprite(width, height, resolution, roadWidth, sprite, scale, destX, destY+bounce, -0.5, -1, 0)
}

func render_fog(x, y, width, height, fog float32) {
	if fog < 1 {
		gfx.SetColor(0, 81, 8, (1-fog)*255)
		gfx.Rect(gfx.FILL, x, y, width, height)
	}
}

func rumbleWidth(projectedRoadWidth, lanes float32) float32 {
	return projectedRoadWidth / float32(math.Max(6, 2*float64(lanes)))
}

func laneMarkerWidth(projectedRoadWidth, lanes float32) float32 {
	return projectedRoadWidth / float32(math.Max(32, 8*float64(lanes)))
}
