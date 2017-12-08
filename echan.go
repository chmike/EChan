/*
Package echan provides a way to insert a buffer between two channels.

Different buffer implementations are available in the subfolders.
*/
package echan

// Interface is the way to insert a buffer between the 'in' and 'out' channels
// When 'in' gets closed:
// 1. All buffered items must be sent to 'out'
// 2. 'out' must be closed
// 3. The function returns
type Interface func(in <-chan interface{}, out chan<- interface{})

// ChanForward synchronously writes all items from in to out
func ChanForward(in <-chan interface{}, out chan<- interface{}) {
	for i := range in {
		out <- i
	}
	close(out)
}
