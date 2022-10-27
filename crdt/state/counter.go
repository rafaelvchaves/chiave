package state

import (
	"fmt"
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

func (c Counter) String() string {
	return c.replica.String() + ": " + c.pos.String() + ", " + c.neg.String()
}

func (c *Counter) Merge(o Counter) {
	fmt.Println(c.String())
	fmt.Println(o.String())
	c.pos.Merge(o.pos)
	c.neg.Merge(o.neg)
}

func (c *Counter) Copy() Counter {
	cpy := NewCounter(c.replica)
	cpy.pos = c.pos.Copy()
	cpy.neg = c.neg.Copy()
	return *cpy
}

func (c *Counter) GetEvent() crdt.Event {
	return crdt.Event{
		Source: c.replica,
		Type:   crdt.CType,
		Data:   c.Copy(),
	}
}

func (s *Counter) PersistEvent(event crdt.Event) {
	c, ok := event.Data.(Counter)
	if !ok {
		return
	}
	s.Merge(c)
}
