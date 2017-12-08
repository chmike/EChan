package echan_test

import (
	"testing"

	"github.com/chmike/EChan"
	"github.com/chmike/EChan/bufferedChannel"
	"github.com/chmike/EChan/immediateWrite"
	"github.com/chmike/EChan/queue"
	"github.com/chmike/EChan/queue/ering"
	"github.com/chmike/EChan/queue/ring"
	"github.com/chmike/EChan/queue/slice"
	etest "github.com/chmike/EChan/testing"
)

func Example() {
	var bc echan.Interface
	// Choose one of the implemenations:

	// ering implements a queue whose capacity my grow and shrink as needed.
	//
	// Under normal usage (output faster than input), memory usage will be minimal.
	// In sporadic congestion conditions, the capacity may grow as needed. When the
	// congestion is resorbed, the internal buffer will shrink and free memory.
	//
	// Nevertheless, an upper capacity limit is defined where input will block in
	// case the output is blocked. This limit is to avoid memory exhaustion and OS hog.
	bc = queue.New(ering.New(5))

	// ring uses a fixed-sized slice as queue backend.
	bc = queue.New(ring.New(5))

	// slice uses a simple slice as queue backend (with an initial capacity).
	// The size is actually unbounded.
	bc = queue.New(slice.New(5))

	// bufferedChannel uses an intermediate buffered channel.
	bc = bufferedChannel.New(5)

	in := make(chan interface{})
	out := make(chan interface{})

	// Producer for 'in'
	go func() {
		for i := 0; i < 50; i++ {
			in <- i
		}
		close(in)
	}()

	// buffering
	go bc(in, out)

	// Consumer for 'ou'
	for _ = range out {
	}
	// Output:
}

func BenchmarkAll(b *testing.B) {
	type factory func(int) echan.Interface
	immediateWriteAdapted := func(int) echan.Interface {
		return immediateWrite.New()
	}
	type qfactory func(int) queue.Interface
	queueAdapted := func(q qfactory) func(int) echan.Interface {
		return func(size int) echan.Interface {
			return queue.New(q(size))
		}
	}
	factories := []struct {
		name string
		f    factory
	}{
		{"bufChan", bufferedChannel.New},
		{"ering", queueAdapted(ering.New)},
		{"ring", queueAdapted(ring.New)},
		{"slice", queueAdapted(slice.New)},
		{"none", immediateWriteAdapted},
	}
	sizes := []struct {
		name     string
		capacity int
		items    int
	}{
		{"Small", 5, 5},
		{"LargeItems-SmallCap", 50, 10000},
	}
	benches := []struct {
		name string
		b    func(*testing.B, echan.Interface, int)
	}{
		{"BufferBoth", etest.BenchmarkBuffBoth},
		{"BufferOut_", etest.BenchmarkBuffOut},
		{"BufferIn__", etest.BenchmarkBuffIn},
		{"BufferNone", etest.BenchmarkBuffNone},
	}

	for _, be := range benches {
		b.Run(be.name, func(b *testing.B) {
			for _, s := range sizes {
				b.Run(s.name, func(b *testing.B) {
					for _, f := range factories {
						b.Run(f.name, func(b *testing.B) {
							bc := f.f(s.capacity)
							be.b(b, bc, s.items)
						})
					}
				})
			}
		})
	}
}

func BenchmarkBuffers(b *testing.B) {
	type factory func(int) echan.Interface
	type qfactory func(int) queue.Interface
	queueAdapted := func(q qfactory) func(int) echan.Interface {
		return func(size int) echan.Interface {
			return queue.New(q(size))
		}
	}
	factories := []struct {
		name string
		f    factory
	}{
		{"bufChan", bufferedChannel.New},
		{"ering", queueAdapted(ering.New)},
		{"ring", queueAdapted(ring.New)},
		{"slice", queueAdapted(slice.New)},
	}
	sizes := []struct {
		name     string
		capacity int
		items    int
	}{
		{"SmallItems-LargeCap", 10000, 50},
		{"LargerItems-LargeCap", 10000, 50000},
	}
	benches := []struct {
		name string
		b    func(*testing.B, echan.Interface, int)
	}{
		{"BufferBoth", etest.BenchmarkBuffBoth},
		{"BufferOut_", etest.BenchmarkBuffOut},
		{"BufferIn__", etest.BenchmarkBuffIn},
		{"BufferNone", etest.BenchmarkBuffNone},
	}

	for _, be := range benches {
		b.Run(be.name, func(b *testing.B) {
			for _, s := range sizes {
				b.Run(s.name, func(b *testing.B) {
					for _, f := range factories {
						b.Run(f.name, func(b *testing.B) {
							bc := f.f(s.capacity)
							be.b(b, bc, s.items)
						})
					}
				})
			}
		})
	}
}
