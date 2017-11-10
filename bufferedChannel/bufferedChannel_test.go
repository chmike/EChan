package bufferedChannel_test

import (
	"testing"

	"github.com/chmike/EChan/bufferedChannel"
	echanTest "github.com/chmike/EChan/testing"
)

func TestAll(t *testing.T) {
	bc := bufferedChannel.New(10)
	echanTest.ImmediateClosing(t, bc)
	echanTest.OneElement(t, bc)
	echanTest.SomeElements(t, bc)
	echanTest.ShouldBlock(t, bc, 12) // +2 because of the goroutines
}

func BenchmarkSimple(b *testing.B) {
	bc := bufferedChannel.New(100)
	echanTest.BenchmarkSimple(b, bc, 100)
}
