package queue

import (
	"github.com/chmike/EChan"
)

// Interface represents a queue which can be used as a backend for buffering items between two channels
type Interface interface {
	Push(interface{}) bool    // return false if queue is full (item was not pushed)
	Pop() (interface{}, bool) // return false if queue is empty
}

// New can store 2 items more than the capacity (due to the way it works)
func New(q Interface) echan.Interface {
	return func(in <-chan interface{}, out chan<- interface{}) {
		var head interface{}
		headOk := false

		defer func() {
			// Empty the queue before closing out
			if headOk {
				out <- head
			}
			head, headOk = q.Pop()
			for headOk {
				out <- head
				head, headOk = q.Pop()
			}
			close(out)
		}()

		for {
			if !headOk {
				head, headOk = q.Pop()
				if !headOk {
					// block for the first element
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
				for !q.Push(new) {
					// queue is full: block until head is out
					out <- head
					head, headOk = q.Pop()
				}
			case out <- head:
				headOk = false
			}
		}
	}
}
