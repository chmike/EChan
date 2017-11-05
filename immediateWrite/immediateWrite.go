package immediateWrite

import (
	"github.com/chmike/EChan"
)

func New() echan.Interface {
	return func(in <-chan interface{}, out chan<- interface{}) {
		for i := range in {
			out <- i
		}
		close(out)
	}
}