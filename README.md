[![GoDoc](https://godoc.org/github.com/chmike/EChan?status.svg)](https://godoc.org/github.com/chmike/EChan)
[![Build](https://travis-ci.org/chmike/EChan.svg?branch=master)](https://travis-ci.org/chmike/EChan?branch=master)
[![Coverage](https://coveralls.io/repos/github/chmike/EChan/badge.svg?branch=master)](https://coveralls.io/github/chmike/EChan?branch=master)
[![Go Report](https://goreportcard.com/badge/github.com/chmike/EChan)](https://goreportcard.com/report/github.com/chmike/EChan)
![Status](https://img.shields.io/badge/status-stable-green.svg)

# EChan, the elastic channel

Package echan provides a channel whose capacity my grow and shrink as needed.

Under normal usage (output faster than input), memory usage will be minimal.
In sporadic congestion conditions, the capacity may grow as needed. When the
congestion is resorbed, the internal buffer will shrink and free memory.

Nevertheless, an upper capacity limit is defined where input will block in
case the output is blocked. This limit is to avoid memory exhaustion and
OS hog.

Example usage:

``` Go
// Instantiating the elastic channel with an upper limit value
c := echan.New(10000)

// Queuing a value in the channel. Blocks if the channel is full.
c.In() <- 123

// Reading a value out of the channel. Blocks while the channel is empty.
v := <-c.Out().(int) 

// Closing the channel
c.Close()
```

Based on an idea of [oliverpool](https://github.com/oliverpool) discussed [here](https://github.com/npat-efault/musings/issues/1#issuecomment-339889714).

**Warning**: The EChan must be closed to be garbage collected. This is because of its internal go routine. If you don't close the EChan, your program will have a memory leak.
