package bufferedChannel

import (
	"github.com/chmike/EChan"
)

// New creates a buffer of the given size (+2) between in and out.
// Two items are additionnaly buffered by the goroutine
func New(size int) echan.Interface {
	return func(in <-chan interface{}, out chan<- interface{}) {
		buf := make(chan interface{}, size)
		go echan.ChanForward(in, buf)
		echan.ChanForward(buf, out)
	}
}
