package slice_test

import (
	"testing"

	"github.com/chmike/EChan/queue"
	"github.com/chmike/EChan/queue/slice"
	echanTest "github.com/chmike/EChan/testing"
)

func TestInterface(t *testing.T) {
	var s queue.Interface
	s = slice.New(5)

	s.Push(42)

	v, _ := s.Pop()
	if v != 42 {
		t.Errorf("Popped should be %d, is %+v", 42, v)
	}

	s.Push(43)
	s.Push(44)
	s.Pop()
	_, ok := s.Pop()
	if !ok {
		t.Error("Could not pop last item")
	}
	_, ok = s.Pop()
	if ok {
		t.Error("Could pop one item to much")
	}
}

func TestEchan(t *testing.T) {
	s := slice.New(5)
	bc := queue.New(s)
	echanTest.ImmediateClosing(t, bc)
	echanTest.OneElement(t, bc)
	echanTest.SomeElements(t, bc)
}

func BenchmarkSimple(b *testing.B) {
	s := slice.New(100)
	bc := queue.New(s)
	echanTest.BenchmarkSimple(b, bc, 100)
}
