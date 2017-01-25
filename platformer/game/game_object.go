package game

import (
	"sort"

	"github.com/tanema/amore-examples/platformer/ump"
)

type (
	gameObject interface {
		update(dt float32)
		tag() string
		destroy()
		push(strength float32)
		damage(intensity float32)
		draw(bool)
		updateOrder() int
	}
	By         func(b1, b2 *ump.Body) bool
	bodySorter struct {
		bodies []*ump.Body
		by     func(b1, b2 *ump.Body) bool
	}
)

func (by By) Sort(bodies []*ump.Body) {
	ps := &bodySorter{
		bodies: bodies,
		by:     by,
	}
	sort.Sort(ps)
}

func (s *bodySorter) Len() int {
	return len(s.bodies)
}

func (s *bodySorter) Swap(i, j int) {
	s.bodies[i], s.bodies[j] = s.bodies[j], s.bodies[i]
}

func (s *bodySorter) Less(i, j int) bool {
	return s.by(s.bodies[i], s.bodies[j])
}
