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

func New(capacity int) queue.CappedInterface {
	return &ring{
		buf:  make([]interface{}, capacity, capacity),
		head: 0,
		len:  0,
		cap:  capacity,
	}
}

func (r ring) Len() int {
	return r.len
}

func (r *ring) Push(e interface{}) {
	r.buf[r.head] = e
	r.head = (r.head + 1) % r.cap
	r.len++
}

func (r *ring) Pop() (x interface{}) {
	x = r.buf[(r.head-r.len+r.cap)%r.cap]
	r.len--
	return x
}

func (r *ring) IsFull() bool {
	return r.len == r.cap
}
