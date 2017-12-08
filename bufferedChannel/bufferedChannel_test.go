package bufferedChannel_test

import (
	"testing"

	"github.com/chmike/EChan/bufferedChannel"
	etest "github.com/chmike/EChan/testing"
)

func TestAll(t *testing.T) {
	bc := bufferedChannel.New(10)
	etest.ImmediateClosing(t, bc)
	etest.OneElement(t, bc)
	etest.SomeElements(t, bc)
	etest.ShouldBlock(t, bc, 12) // +2 because of the goroutines
}

func BenchmarkBuffBoth(b *testing.B) {
	bc := bufferedChannel.New(100)
	etest.BenchmarkBuffBoth(b, bc, 100)
}
