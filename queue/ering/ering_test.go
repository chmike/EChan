package ering_test

import (
	"testing"

	"github.com/chmike/EChan/queue"
	"github.com/chmike/EChan/queue/ering"
	echanTest "github.com/chmike/EChan/testing"
)

func TestInterface(t *testing.T) {
	var s queue.Interface
	s = ering.New(2)

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

	s = ering.New(100)
	for i := 0; i < 100; i++ {
		if !s.Push(i) {
			t.Errorf("Could not push %d", i)
		}
	}
	for i := 0; i < 100; i++ {
		v, ok := s.Pop()
		if !ok {
			t.Errorf("Pop returned false instead of %d", i)
		} else if v != i {
			t.Errorf("got %v, expected %d", v, i)
		}
	}

	for i := 0; i < 30; i++ {
		if !s.Push(i) {
			t.Errorf("Could not push %d", i)
		}
	}
	// make ring buffer wrap
	for i := 0; i < 10; i++ {
		v, ok := s.Pop()
		if !ok {
			t.Errorf("Pop returned false instead of %d", i)
		} else if v != i {
			t.Errorf("got %v, expected %d", v, i)
		}
		if !s.Push(i + 30) {
			t.Errorf("Could not push %d", i+30)
		}
	}
	// empty buffer and make it shrink to minimal size
	for i := 10; i < 40; i++ {
		v, ok := s.Pop()
		if !ok {
			t.Errorf("Pop returned false instead of %d", i)
		} else if v != i {
			t.Errorf("got %v, expected %d", v, i)
		}
	}
}

func TestEchanCapped(t *testing.T) {
	s := ering.New(5)
	bc := queue.New(s)
	echanTest.ImmediateClosing(t, bc)
	echanTest.OneElement(t, bc)
	echanTest.SomeElements(t, bc)
	echanTest.ShouldBlock(t, bc, 7) // +2 because of the goroutines
}

func BenchmarkSimple(b *testing.B) {
	s := ering.New(100)
	bc := queue.New(s)
	echanTest.BenchmarkSimple(b, bc, 100)
}
