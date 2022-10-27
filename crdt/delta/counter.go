package delta

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

type delta = struct {
	posDelta map[string]int
	negDelta map[string]int
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

func (c *Counter) Merge(d delta) {
	c.pos.Merge(d.posDelta)
	c.neg.Merge(d.negDelta)
}

func (c *Counter) String() string {
	return fmt.Sprintf("%d", c.Value())
}

func (c *Counter) GetEvent() crdt.Event {
	return crdt.Event{
		Source: c.replica,
		Type:   crdt.CType,
		Data: delta{
			posDelta: c.pos.GetDelta(),
			negDelta: c.neg.GetDelta(),
		},
	}
}

func (s *Counter) PersistEvent(event crdt.Event) {
	d, ok := event.Data.(delta)
	if !ok {
		return
	}
	s.Merge(d)
}
