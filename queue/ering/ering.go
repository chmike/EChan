/*
Package ering implements a queue whose capacity my grow and shrink as needed.

Under normal usage (output faster than input), memory usage will be minimal.
In sporadic congestion conditions, the capacity may grow as needed. When the
congestion is resorbed, the internal buffer will shrink and free memory.

Nevertheless, an upper capacity limit is defined where input will block in
case the output is blocked. This limit is to avoid memory exhaustion and OS hog.
*/
package ering

import (
	"github.com/chmike/EChan/queue"
)

type ering struct {
	buf  []interface{}
	head int
	len  int
	cap  int
}

// New instatiate a ring queue whose internal buffer will grow and shrink as needed.
// Capacity is the maximum number of elements that can be stored in the queue.
func New(capacity int) queue.Interface {
	return &ering{
		buf:  make([]interface{}, 4),
		head: 0,
		len:  0,
		cap:  capacity,
	}
}

func (r *ering) Push(e interface{}) bool {
	if r.len == r.cap {
		return false
	}
	if r.len == len(r.buf) {
		// grow buffer
		var tmp = make([]interface{}, 2*len(r.buf))
		copy(tmp[copy(tmp, r.buf[r.head:]):], r.buf[:r.head])
		r.buf = tmp
		r.head = r.len
	}
	r.buf[r.head] = e
	r.head = (r.head + 1) % len(r.buf)
	r.len++
	return true
}

func (r *ering) Pop() (x interface{}, ok bool) {
	if r.len == 0 {
		return nil, false
	}
	var pos = (r.head - r.len + len(r.buf)) % len(r.buf)
	x = r.buf[pos]
	r.buf[pos] = nil // to avoid memory leaks
	r.len--
	if r.len > 1 && r.len == len(r.buf)/4 {
		// shrink buffer
		var tmp = make([]interface{}, len(r.buf)/2)
		if r.len < r.head {
			copy(tmp, r.buf[r.head-r.len:r.head])
		} else {
			copy(tmp[copy(tmp, r.buf[len(r.buf)-(r.len-r.head):]):], r.buf[:r.head])
		}
		r.buf = tmp
		r.head = r.len
	}
	return x, true
}
