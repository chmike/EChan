package slice

import (
	"github.com/chmike/EChan/queue"
)

type slice struct {
	s []interface{}
}

type cappedSlice struct {
	slice
	cap int
}

func New(capacity int) queue.Interface {
	return &slice{
		s: make([]interface{}, 0, capacity),
	}
}

func NewCapped(capacity int) queue.CappedInterface {
	return &cappedSlice{
		slice: slice{make([]interface{}, 0, capacity)},
		cap:   capacity,
	}
}

func (s slice) Len() int {
	return len(s.s)
}

func (s *slice) Pop() (x interface{}) {
	x, s.s = s.s[0], s.s[1:]
	return x
}

func (s *slice) Push(e interface{}) {
	s.s = append(s.s, e)
}

func (s *cappedSlice) IsFull() bool {
	return len(s.s) == s.cap
}
