package game

import (
	"math"
)

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
