package bufferedChannel

import (
	"github.com/chmike/EChan"
)

func New(size int) echan.Interface {
	if size >= 2 {
		// Two are already "buffered" via the goroutine
		size -= 2
	}
	return func(in <-chan interface{}, out chan<- interface{}) {
		buf := make(chan interface{}, size)
		go echan.ChanForward(in, buf)
		echan.ChanForward(buf, out)
	}
}
