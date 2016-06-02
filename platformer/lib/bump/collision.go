package bump

type (
	Point struct {
		X, Y float32
	}
	Collision struct {
		Body     *Body
		Other    *Body
		Move     Point
		Normal   Point
		Touch    Point
		Data     Point
		RespType string
		overlaps bool
		ti       float32
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
	if a.ti == b.ti {
		return a.Body.distanceTo(a.Other) < b.Body.distanceTo(b.Other)
	}
	return a.ti < b.ti
}
