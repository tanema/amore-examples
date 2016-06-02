package game

import (
	"math/rand"
)

func randMax(max float32) float32 {
	return randRange(0, max)
}

func randLimits(limit float32) float32 {
	return randRange(-limit, limit)
}

func randRange(min, max float32) float32 {
	return (rand.Float32() * (max - min)) + min
}
