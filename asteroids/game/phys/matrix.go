package phys

import (
	"math"
)

type Matrix [6]float32

// rot in radians
func (mat *Matrix) translate(x, y, rot, scale float32) {
	sin := float32(math.Sin(float64(rot))) * scale
	cos := float32(math.Cos(float64(rot))) * scale
	mat[0], mat[1], mat[2], mat[3], mat[4], mat[5] = cos, -sin, x, sin, cos, y
}

func (mat *Matrix) multiply(x, y float32) (float32, float32) {
	return (mat[0] * x) + (mat[1] * y) + mat[2], (mat[3] * x) + (mat[4] * y) + mat[5]
}
