/*
Package echan provides a channel whose capacity my grow and shrink as needed.

Under normal usage (output faster than input), memory usage will be minimal.
In sporadic congestion conditions, the capacity may grow as needed. When the
congestion is resorbed, the internal buffer will shrink and free memory.

Nevertheless, an upper capacity limit is defined where input will block in
case the output is blocked. This limit is to avoid memory exhaustion and
OS hog.
*/
package echan

type Implementation func(in <-chan interface{}, out chan<- interface{})

const minBufCap = 8
const chanCap = 8

// EChan is a channel whose internal queue capacity my shrink and grow as needed.
type EChan struct {
	in, out            chan interface{}
	done               chan struct{}
	buf                []interface{}
	beg, end, cnt, max int
}

// New returns a new file channel.
func New(max int) *EChan {
	max -= 2 * chanCap
	if max < 2*minBufCap {
		max = 2 * minBufCap
	}
	var c = &EChan{
		in:   make(chan interface{}, chanCap),
		out:  make(chan interface{}, chanCap),
		done: make(chan struct{}),
		buf:  make([]interface{}, minBufCap),
		max:  max,
	}
	go c.run()
	return c
}

// Close close the channel.
func (c *EChan) Close() {
	close(c.done)
	close(c.in)
}

// In returns the input channel. Don't close it. Use the Close() method.
func (c *EChan) In() chan interface{} {
	return c.in
}

// Out returns the output channel. Don't close it. Use the Close() method.
func (c *EChan) Out() chan interface{} {
	return c.out
}

func (c *EChan) pushBack(v interface{}) {
	if c.cnt == len(c.buf) {
		var tmp = make([]interface{}, 2*c.cnt)
		copy(tmp[copy(tmp, c.buf[c.beg:]):], c.buf[:c.end])
		c.beg = 0
		c.end = c.cnt
		c.buf = tmp
	}
	c.buf[c.end] = v
	c.cnt++
	c.end++
	if c.end == len(c.buf) {
		c.end = 0
	}
}

func (c *EChan) popFront() interface{} {
	if c.cnt == 0 {
		return nil
	}
	var v = c.buf[c.beg]
	c.buf[c.beg] = nil
	c.cnt--
	c.beg++
	if c.beg == len(c.buf) {
		c.beg = 0
	}
	if c.cnt <= len(c.buf)/4 && len(c.buf) >= (minBufCap*4) {
		var tmp = make([]interface{}, len(c.buf)/2)
		if c.beg < c.end {
			copy(tmp, c.buf[c.beg:c.end])
		} else {
			copy(tmp[copy(tmp, c.buf[c.beg:]):], c.buf[:c.end])
		}
		c.beg = 0
		c.end = c.cnt
		c.buf = tmp
	}
	return v
}

func (c *EChan) run() {
	var vi, vo interface{}
	var ok bool
start:
	vo, ok = <-c.in
	for ok {
		select {
		case vi, ok = <-c.in:
			if ok {
				c.pushBack(vi)
				for c.cnt > c.max {
					vo = c.popFront()
					select {
					case c.out <- vo:
					case <-c.done:
						ok = false
					}
				}
			}
		case c.out <- vo:
			if vo = c.popFront(); vo == nil {
				goto start
			}
		}
	}
	close(c.out)
	c.in = nil
	c.out = nil
	c.done = nil
	c.buf = nil
}
