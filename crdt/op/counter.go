package op

import "kvs/crdt"

type OCounter struct {
	c int
}

var _ crdt.Counter = NewOCounter()

func NewOCounter() *OCounter {
	return &OCounter{
		c: 0,
	}
}

func (o *OCounter) Value() int {
	return int(o.c)
}

func (o *OCounter) Increment() {
	o.c += 1
}

func (o *OCounter) Decrement() {
	o.c -= 1
}
