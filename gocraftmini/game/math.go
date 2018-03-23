package game

import (
	"math"
	"math/rand"
	"time"
)

func max(x, y float32) float32 {
	return float32(math.Max(float64(x), float64(y)))
}

func pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

func abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func exp(x float32) float32 {
	return float32(math.Exp(float64(x)))
}

func sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

func randRange(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min+1) + min
}
