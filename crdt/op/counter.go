package op

import (
	"fmt"
	"kvs/crdt"
	"kvs/util"
)

type Counter struct {
	replica util.Replica
	c       int
	current crdt.Event
}

func NewCounter(replica util.Replica) *Counter {
	return &Counter{
		replica: replica,
		c:       0,
		current: crdt.Event{
			Source: replica,
		},
	}
}

func (c *Counter) Value() int {
	return int(c.c)
}

func (c *Counter) String() string {
	return fmt.Sprintf("%d", c.Value())
}

func (c *Counter) Increment() {
	c.c += 1
	update, _ := c.current.Data.(int)
	c.current.Data = update + 1
}

func (c *Counter) Decrement() {
	c.c -= 1
	update, _ := c.current.Data.(int)
	c.current.Data = update - 1
}

func (c *Counter) GetEvent() crdt.Event {
	current := c.current
	c.current = crdt.Event{
		Source: c.replica,
		Data:   0,
	}
	return current
}

func (c *Counter) PersistEvent(event crdt.Event) {
	update, ok := event.Data.(int)
	if !ok {
		return
	}
	c.c += update
}
