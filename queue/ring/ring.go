package ring

import (
	"github.com/chmike/EChan/queue"
)

type ring struct {
	buf  []interface{}
	head int
	len  int
	cap  int
}

func New(capacity int) queue.Interface {
	return &ring{
		buf:  make([]interface{}, capacity, capacity),
		head: 0,
		len:  0,
		cap:  capacity,
	}
}

func (r *ring) Push(e interface{}) bool {
	if r.len == r.cap {
		return false
	}
	r.buf[r.head] = e
	r.head = (r.head + 1) % r.cap
	r.len++
	return true
}

func (r *ring) Pop() (x interface{}, ok bool) {
	if r.len == 0 {
		return nil, false
	}
	x = r.buf[(r.head-r.len+r.cap)%r.cap]
	r.len--
	return x, true
}
