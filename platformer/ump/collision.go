package ump

type (
	Point struct {
		X, Y float32
	}
	Collision struct {
		Intersection float32
		Distance     float32
		Body         *Body
		Move         Point
		Normal       Point
		Touch        Point
		Data         Point
		RespType     string
	}
)

type CollisionsByDistance []*Collision

func (s CollisionsByDistance) Len() int {
	return len(s)
}

func (s CollisionsByDistance) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s CollisionsByDistance) Less(i, j int) bool {
	a, b := s[i], s[j]
	if a.Intersection == b.Intersection {
		return a.Distance < b.Distance
	}
	return a.Intersection < b.Intersection
}
