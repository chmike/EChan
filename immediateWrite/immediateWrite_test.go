package immediateWrite_test

import (
	"testing"

	"github.com/chmike/EChan/immediateWrite"
	echanTest "github.com/chmike/EChan/testing"
)

func TestAll(t *testing.T) {
	bc := immediateWrite.New()
	echanTest.ImmediateClosing(t, bc)
	echanTest.OneElement(t, bc)
	echanTest.SomeElements(t, bc)
	echanTest.ShouldBlock(t, bc, 1)
}

func BenchmarkSimple(b *testing.B) {
	bc := immediateWrite.New()
	echanTest.BenchmarkSimple(b, bc, 100)
}
