package state

import "kvs/crdt"

type Counter struct {
	id  string
	pos GCounter
	neg GCounter
}

func NewCounter(id string) *Counter {
	return &Counter{
		id:  id,
		pos: NewGCounter(id),
		neg: NewGCounter(id),
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

func (c *Counter) GetEvent() crdt.Event {
	return crdt.Event{
		Source: c.id,
		Data:   *c,
	}
}

func (s *Counter) PersistEvent(event crdt.Event) {
	s.Merge(event.Data.(Counter))
}
