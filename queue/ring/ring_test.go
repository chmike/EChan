package ring_test

import (
	"testing"

	"github.com/chmike/EChan/queue"
	"github.com/chmike/EChan/queue/ring"
	echanTest "github.com/chmike/EChan/testing"
)

func TestInterface(t *testing.T) {
	var s queue.Interface
	s = ring.New(2)

	s.Push(42)

	v, _ := s.Pop()
	if v != 42 {
		t.Errorf("Popped should be %d, is %+v", 42, v)
	}

	s.Push(43)
	s.Push(44)
	if s.Push(45) {
		t.Error("Could push too many items")
	}
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

func TestEchanCapped(t *testing.T) {
	s := ring.New(5)
	bc := queue.New(s)
	echanTest.ImmediateClosing(t, bc)
	echanTest.OneElement(t, bc)
	echanTest.SomeElements(t, bc)
	echanTest.ShouldBlock(t, bc, 7) // +2 because of the goroutines
}

func BenchmarkSimple(b *testing.B) {
	s := ring.New(100)
	bc := queue.New(s)
	echanTest.BenchmarkSimple(b, bc, 100)
}
