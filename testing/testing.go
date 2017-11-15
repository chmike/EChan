package testing

import (
	"testing"
	"time"

	echan "github.com/chmike/EChan"
)

func ImmediateClosing(t *testing.T, imp echan.Interface) {
	in := make(chan interface{})
	out := make(chan interface{})
	go close(in)
	imp(in, out)
}

func OneElement(t *testing.T, imp echan.Interface) {
	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		in <- 42
		close(in)
	}()
	go imp(in, out)
	for o := range out {
		v, ok := o.(int)
		if !ok {
			t.Errorf("Expected int, got %T", o)
		} else if v != 42 {
			t.Errorf("Expected 42, got %#v", v)
		}
	}
}

func SomeElements(t *testing.T, imp echan.Interface) {
	in := make(chan interface{})
	out := make(chan interface{})
	elts := []interface{}{
		1, "lol", nil, 1 * time.Second,
	}
	go func() {
		for _, e := range elts {
			in <- e
		}
		close(in)
	}()
	go imp(in, out)
	i := 0
	for o := range out {
		if i >= len(elts) {
			t.Error("Got more elements than expected")
			return
		}
		expected := elts[i]
		i++
		if expected != o {
			t.Errorf("Expected %+v, got %+v", expected, o)
		}
	}
}

func ShouldBlock(t *testing.T, imp echan.Interface, size int) {
	in := make(chan interface{})
	out := make(chan interface{})
	sigOut := make(chan struct{})
	go func() {
		for i := 0; i < size; i++ {
			t.Logf("Send %d", i)
			in <- i
		}
		select {
		case in <- -1:
			t.Logf("Sent %d", -1)
			t.Errorf("Shouldn't have accepted a new element (after %d send)", size)
			close(sigOut)
		case <-time.After(1 * time.Millisecond):
			t.Log("Blocked succesfully")
			close(sigOut)
			t.Logf("Send %d", -2)
			in <- -2
		}
		close(in)
	}()
	go imp(in, out)

	<-sigOut
	exp := 0
	for o := range out {
		v, ok := o.(int)
		if !ok {
			t.Errorf("Expected int, got %T", o)
		} else if v != exp {
			t.Errorf("Expected %d, got %d", exp, v)
		} else {
			t.Logf("Received %d", v)
		}
		exp++
		if exp == size {
			exp = -2
		}
	}
}

func BenchmarkBuffBoth(b *testing.B, imp echan.Interface, size int) {
	for n := 0; n < b.N; n++ {
		in := make(chan interface{}, size)
		out := make(chan interface{}, size)
		for i := 0; i < size; i++ {
			in <- i
		}
		close(in)
		imp(in, out)
		for _ = range out {
		}
	}
}

func BenchmarkBuffIn(b *testing.B, imp echan.Interface, size int) {
	for n := 0; n < b.N; n++ {
		in := make(chan interface{}, size)
		out := make(chan interface{})
		for i := 0; i < size; i++ {
			in <- i
		}
		close(in)
		go imp(in, out)
		for _ = range out {
		}
	}
}

func BenchmarkBuffOut(b *testing.B, imp echan.Interface, size int) {
	for n := 0; n < b.N; n++ {
		in := make(chan interface{})
		out := make(chan interface{}, size)
		go func() {
			for i := 0; i < size; i++ {
				in <- i
			}
			close(in)
		}()
		imp(in, out)
		for _ = range out {
		}
	}
}
