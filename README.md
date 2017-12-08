[![GoDoc](https://godoc.org/github.com/chmike/EChan?status.svg)](https://godoc.org/github.com/chmike/EChan)
[![Build](https://travis-ci.org/chmike/EChan.svg?branch=master)](https://travis-ci.org/chmike/EChan?branch=master)
[![Coverage](https://coveralls.io/repos/github/chmike/EChan/badge.svg?branch=master)](https://coveralls.io/github/chmike/EChan?branch=master)
[![Go Report](https://goreportcard.com/badge/github.com/chmike/EChan)](https://goreportcard.com/report/github.com/chmike/EChan)
![Status](https://img.shields.io/badge/status-stable-brightgreen.svg)

# EChan, the elastic channel

Package echan provides a way to insert a buffer between two channels.

You can switch the implementation depending on your scenario.

Example usage:

``` Go
package main

import (
	"github.com/chmike/EChan"
	"github.com/chmike/EChan/bufferedChannel"
	"github.com/chmike/EChan/queue"
	"github.com/chmike/EChan/queue/ering"
	"github.com/chmike/EChan/queue/ring"
	"github.com/chmike/EChan/queue/slice"
)

func main() {
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
}

```

Based on an idea of [oliverpool](https://github.com/oliverpool) discussed [here](https://github.com/npat-efault/musings/issues/1#issuecomment-339889714).
