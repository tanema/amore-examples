package ump

import (
	"sort"
)

type (
	By         func(b1, b2 *Body) bool
	bodySorter struct {
		bodies []*Body
		by     func(b1, b2 *Body) bool
	}
)

func (by By) Sort(bodies []*Body) {
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
