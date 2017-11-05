package queue

import (
	"io"

	"github.com/chmike/EChan"
)

type Interface interface {
	Len() int
	Push(interface{})
	Pop() interface{}
}

type CappedInterface interface {
	Interface
	IsFull() bool
}

func New(q Interface) echan.Interface {
	return func(in <-chan interface{}, out chan<- interface{}) {
		var head interface{}
		headOk := false

		defer func() {
			// Empty the queue before closing out
			if headOk {
				out <- head
			}
			for q.Len() > 0 {
				out <- q.Pop()
			}
			close(out)
			if c, ok := q.(io.Closer); ok {
				c.Close()
			}
		}()

		for {
			if !headOk {
				if q.Len() > 0 {
					head = q.Pop()
					headOk = true
				} else {
					head, headOk = <-in
					if !headOk {
						// in channel closed
						return
					}
				}
			}
			// we have a valid head
			select {
			case new, newOk := <-in:
				if !newOk {
					// in channel closed
					return
				}
				q.Push(new)
			case out <- head:
				headOk = false
			}
		}
	}
}

// NewCapped can store 2 more than the capacity (due to the way it works)
func NewCapped(q CappedInterface) echan.Interface {
	return func(in <-chan interface{}, out chan<- interface{}) {
		var head interface{}
		headOk := false

		defer func() {
			// Empty the queue before closing out
			if headOk {
				out <- head
			}
			for q.Len() > 0 {
				out <- q.Pop()
			}
			close(out)
			if c, ok := q.(io.Closer); ok {
				c.Close()
			}
		}()

		for {
			if !headOk {
				if q.Len() > 0 {
					head = q.Pop()
					headOk = true
				} else {
					head, headOk = <-in
					if !headOk {
						// in channel closed
						return
					}
				}
			}
			// we have a valid head
			select {
			case new, newOk := <-in:
				if !newOk {
					// in channel closed
					return
				}
				// queue is full: block until head is out
				for q.IsFull() {
					out <- head
					head = q.Pop()
				}
				q.Push(new)
			case out <- head:
				headOk = false
			}
		}
	}
}
