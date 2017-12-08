package immediateWrite_test

import (
	"testing"

	"github.com/chmike/EChan/immediateWrite"
	etest "github.com/chmike/EChan/testing"
)

func TestAll(t *testing.T) {
	bc := immediateWrite.New()
	etest.ImmediateClosing(t, bc)
	etest.OneElement(t, bc)
	etest.SomeElements(t, bc)
	etest.ShouldBlock(t, bc, 1)
}

func BenchmarkBuffBoth(b *testing.B) {
	bc := immediateWrite.New()
	etest.BenchmarkBuffBoth(b, bc, 100)
}
