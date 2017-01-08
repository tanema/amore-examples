package ump

import (
	"math"
	"math/rand"
)

var (
	inf           = float32(math.Inf(1))
	delta float32 = 1e-10
)

func clamp(x, minX, maxX float32) float32 {
	if x < minX {
		return minX
	} else if x > maxX {
		return maxX
	}
	return x
}

func randMax(max float32) float32 {
	return randRange(0, max)
}

func randRange(min, max float32) float32 {
	return (rand.Float32() * (max - min)) + min
}

func randLimits(limit float32) float32 {
	return randRange(-limit, limit)
}

func abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func min(x, y float32) float32 {
	return float32(math.Min(float64(x), float64(y)))
}

func max(x, y float32) float32 {
	return float32(math.Max(float64(x), float64(y)))
}

func sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func sign(x float32) float32 {
	if x > 0 {
		return 1
	}
	if x == 0 {
		return 0
	}
	return -1
}

func crossProduct(x1, y1, x2, y2 float32) float32 {
	return x1*y2 - y1*x2
}

func nearest(x, a, b float32) float32 {
	if abs(a-x) < abs(b-x) {
		return a
	}
	return b
}

func frac(n float32) float32 {
	return abs(n - float32(math.Floor(float64(n))))
}
