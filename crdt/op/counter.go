package op

import "kvs/crdt"

type Counter struct {
	id      string
	key     string
	c       int
	current crdt.Event
}

func NewCounter(id string, key string) *Counter {
	return &Counter{
		id:  id,
		key: key,
		c:   0,
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
		Data: 0,
	}
	return current
}

func (c *Counter) PersistEvents(events []crdt.Event) {
	for _, e := range events {
		update, ok := e.Data.(int)
		if !ok {
			continue
		}
		c.c += update
	}
}
