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
	echanTest.ShouldBlock(t, bc, 10)
}
