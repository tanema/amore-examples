package game

import (
	"github.com/tanema/amore/gfx"
)

type Car struct {
	*Sprite
	percent float32
	point   Point
	speed   float32
	segment *Segment
}

func newCar(segment *Segment, source *gfx.Quad, z, offset, speed float32) *Car {
	new_car := &Car{
		point: Point{
			z: z,
			w: source.GetWidth() * sprite_scale,
		},
		Sprite:  newSprite(source, offset),
		speed:   speed,
		segment: segment,
	}
	return new_car
}

func (car *Car) update(dt float32, player *Player) {
	car.steer(player)
	car.point.z = increase(car.point.z, dt*car.speed, trackLength)
	car.percent = percentRemaining(car.point.z, segmentLength) // useful for interpolation during rendering phase
	newSegment := findSegment(car.point.z)
	if car.segment != newSegment {
		for i, c := range car.segment.cars {
			if car == c {
				car.segment.cars = append(car.segment.cars[:i], car.segment.cars[i+1:]...)
				break
			}
		}
		newSegment.cars = append(newSegment.cars, car)
		car.segment = newSegment
	}
}

func (car *Car) steer(player *Player) {
	var dir float32
	lookahead := 20

	// optimization, dont bother steering around other cars when 'out of sight' of the player
	if float32(car.segment.index-player.segment.index) > drawDistance {
		return
	}

	for i := 1; i < lookahead; i++ {
		segment := segments[(car.segment.index+i)%len(segments)]

		if (segment == player.segment) && (car.speed > player.speed) && (overlap(player.point.x, player.point.w, car.offset, car.point.w, 1.2)) {
			if player.point.x > 0.5 {
				dir = -1
			} else if (player.point.x < -0.5) || car.offset > player.point.x {
				dir = 1
			} else {
				dir = -1
			}
			car.offset += dir * 1 / float32(i) * (car.speed - player.speed) / maxSpeed // the closer the cars (smaller i) and the greated the speed ratio, the larger the offset
			return
		}

		for j := 0; j < len(segment.cars); j++ {
			otherCar := segment.cars[j]
			if (car.speed > otherCar.speed) && overlap(car.offset, car.point.w, otherCar.offset, otherCar.point.w, 1.2) {
				if otherCar.offset > 0.5 {
					dir = -1
				} else if (otherCar.offset < -0.5) || (car.offset > otherCar.offset) {
					dir = 1
				} else {
					dir = -1
				}
				car.offset += dir * 1 / float32(i) * (car.speed - otherCar.speed) / maxSpeed
				return
			}
		}
	}

	// if no cars ahead, but I have somehow ended up off road, then steer back on
	if car.offset < -0.9 {
		car.offset += 0.1
	} else if car.offset > 0.9 {
		car.offset += -0.1
	}
}
