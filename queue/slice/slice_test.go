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

	assertLen(t, s, 0)

	s.Push(42)
	assertLen(t, s, 1)

	v := s.Pop()
	if v != 42 {
		t.Errorf("Popped should be %d, is %+v", 42, v)
	}
	assertLen(t, s, 0)

	s.Push(43)
	s.Push(44)
	assertLen(t, s, 2)
	s.Pop()
	s.Pop()
	assertLen(t, s, 0)
}

func assertLen(t *testing.T, q queue.Interface, expected int) {
	l := q.Len()
	if l != expected {
		t.Errorf("Len should be %d, is %d", expected, l)
	}
}

func TestEchan(t *testing.T) {
	s := slice.New(5)
	bc := queue.New(s)
	echanTest.ImmediateClosing(t, bc)
	echanTest.OneElement(t, bc)
	echanTest.SomeElements(t, bc)
	// echanTest.ShouldBlock(t, bc, 1)
}

func TestEchanCapped(t *testing.T) {
	s := slice.NewCapped(5)
	bc := queue.NewCapped(s)
	echanTest.ImmediateClosing(t, bc)
	echanTest.OneElement(t, bc)
	echanTest.SomeElements(t, bc)
	echanTest.ShouldBlock(t, bc, 7)
}
