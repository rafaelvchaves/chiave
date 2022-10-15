package state

import (
	"kvs/crdt"
	"kvs/util"
)

type Counter struct {
	replica util.Replica
	pos     GCounter
	neg     GCounter
}

func NewCounter(replica util.Replica) *Counter {
	return &Counter{
		replica: replica,
		pos:     NewGCounter(replica),
		neg:     NewGCounter(replica),
	}
}

func (c *Counter) Value() int {
	return c.pos.Value() - c.neg.Value()
}

func (c *Counter) Increment() {
	c.pos.Increment()
}

func (c *Counter) Decrement() {
	c.neg.Increment()
}

func (c *Counter) Merge(o Counter) {
	c.pos.Merge(o.pos)
	c.neg.Merge(o.neg)
}

func (c *Counter) GetEvent() crdt.Event[CRDT] {
	return crdt.Event[CRDT]{
		Source: c.replica,
		Data:   *c,
	}
}

func (s *Counter) PersistEvent(event crdt.Event[CRDT]) {
	c, ok := event.Data.(Counter)
	if !ok {
		return
	}
	s.Merge(c)
}
