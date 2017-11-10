package slice

import (
	"github.com/chmike/EChan/queue"
)

type slice struct {
	s []interface{}
}

func New(capacity int) queue.Interface {
	return &slice{
		s: make([]interface{}, 0, capacity),
	}
}

func (s *slice) Pop() (x interface{}, ok bool) {
	if len(s.s) == 0 {
		return nil, false
	}
	x, s.s = s.s[0], s.s[1:]
	return x, true
}

func (s *slice) Push(e interface{}) bool {
	s.s = append(s.s, e)
	return true
}
