package bufferedChannel

import (
	"github.com/chmike/EChan"
)

func New(size int) echan.Implementation {
	if size >= 2 {
		// Two are already "buffered" via the goroutine
		size -= 2
	}
	return func(in <-chan interface{}, out chan<- interface{}) {
		buf := make(chan interface{}, size)
		go func(b chan<- interface{}) {
			for i := range in {
				b <- i
			}
			close(b)
		}(buf)

		for b := range buf {
			out <- b
		}
		close(out)
	}
}
