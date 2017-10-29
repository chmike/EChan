package echan

import (
	"sync"
	"testing"
	"time"
)

func TestEChan1(t *testing.T) {
	var c = New(200)
	for i := 0; i < 100; i++ {
		c.In() <- i
	}
	for i := 0; i < 100; i++ {
		<-c.Out()
	}
	c.Close()
	time.Sleep(100 * time.Millisecond)
	if c.buf != nil {
		t.Error("expected channel to be closed")
	}
}

func TestEChan2(t *testing.T) {
	var c = New(1000)
	var in = 1000000
	var out int
	var w sync.WaitGroup
	w.Add(1)
	go func() {
		for out < in {
			if _, ok := <-c.Out(); !ok {
				break
			}
			out++
		}
		w.Done()
	}()
	for i := 0; i < in; i++ {
		c.In() <- i
	}
	w.Wait()
	if out != in {
		t.Errorf("got %d, expected %d", out, in)
	}
	c.Close()
}

func TestEmptyQueue(t *testing.T) {
	var c = EChan{buf: make([]interface{}, 2)}

	if c.popFront() != nil {
		t.Error("expected nil file")
	}
	c.pushBack(nil)
	if c.popFront() != nil {
		t.Error("expected nil file")
	}

	c.pushBack(1)
	c.pushBack(2)
	var fo = c.popFront()
	if fo.(int) != 1 {
		t.Errorf("got %d, expected 1", fo.(int))
	}
	fo = c.popFront()
	if fo.(int) != 2 {
		t.Errorf("got %d, expected 2", fo.(int))
	}
	if c.popFront() != nil {
		t.Error("expected nil file")
	}
}

func TestGrowQueue(t *testing.T) {
	var c = EChan{buf: make([]interface{}, 2)}

	c.pushBack(1)
	c.pushBack(2)
	c.pushBack(3)
	if len(c.buf) != 4 {
		t.Errorf("got %d, expected 4", len(c.buf))
	}

	c.popFront()
	c.pushBack(4)
	c.pushBack(5)
	c.pushBack(6)
	if len(c.buf) != 8 {
		t.Errorf("got %d, expected 8", len(c.buf))
	}

	for i := 2; i <= 6; i++ {
		if f := c.popFront(); f == nil || f.(int) != i {
			if f == nil {
				t.Error("unexpected nil file")
			} else {
				t.Errorf("got %d, expected %d", f.(int), i)
			}
		}
	}
	if c.popFront() != nil {
		t.Error("expected nil file")
	}

}

func TestShrinkQueue(t *testing.T) {
	var c = EChan{buf: make([]interface{}, 2)}
	for i := 0; i < 256; i++ {
		c.pushBack(i)
	}
	for i := 0; i < 256; i++ {
		if f := c.popFront(); f == nil || f.(int) != i {
			if f == nil {
				t.Error("unexpected nil file")
			} else {
				t.Errorf("got %d, expected %d", f.(int), i)
			}
		}
	}
	if c.popFront() != nil {
		t.Error("expected nil file")
	}
	if len(c.buf) != 2*minBufCap {
		t.Errorf("got len %d, expected %d", len(c.buf), 2*minBufCap)
	}
}

func TestMaxCapacity(t *testing.T) {
	var c = New(0)

	if c.max != 2*minBufCap {
		t.Errorf("got %d, expected %d", c.max, 2*minBufCap)
	}

	c = New(100)
	if c.max != 100-2*chanCap {
		t.Errorf("got %d, expected %d", c.max, 100-2*chanCap)
	}
}

func TestCloseEChan1(t *testing.T) {
	var c = New(100)
	time.Sleep(100 * time.Millisecond)
	c.Close()
	time.Sleep(100 * time.Millisecond)
	if c.buf != nil {
		t.Error("expected nil buf")
	}
}

func TestCloseEChan2(t *testing.T) {
	var c = New(1000)

	time.Sleep(100 * time.Millisecond)
	for i := 0; i < 200; i++ {
		c.In() <- i
	}
	c.Close()
	time.Sleep(100 * time.Millisecond)
	if c.buf != nil {
		t.Error("expected nil buf")
	}
}

func TestCloseEChan3(t *testing.T) {
	var c = New(100)
	go func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic error")
			}
		}()
		for i := 0; i < 200; i++ {
			c.In() <- i
		}
		t.Error("expected a panic")
	}()
	time.Sleep(250 * time.Millisecond)
	c.Close()
	time.Sleep(100 * time.Millisecond)
	if c.buf != nil {
		t.Error("expected nil buf")
	}
}
