package op

import "kvs/crdt"

type Counter struct {
	id      string
	c       int
	current crdt.Event
}

func NewCounter(id string) *Counter {
	return &Counter{
		id: id,
		c:  0,
		current: crdt.Event{
			Source: id,
		},
	}
}

func (c *Counter) Value() int {
	return int(c.c)
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
		Source: c.id,
	}
	return current
}

func (c *Counter) PersistEvent(event crdt.Event) {
	update, _ := event.Data.(int)
	c.c += update
}
