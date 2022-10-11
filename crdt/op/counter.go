package op

import "kvs/crdt"

type OCounter struct {
	c int
	events []crdt.Event
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
	o.events = append(o.events, crdt.Event{Data: 1})
}

func (o *OCounter) Decrement() {
	o.c -= 1
	o.events = append(o.events, crdt.Event{Data: -1})
}

func (o *OCounter) GetEvents() []crdt.Event {
	events := o.events
	o.events = nil
	return events
}

func (o *OCounter) PersistEvents(events []crdt.Event) {
	for _, e := range events {
		if d, ok := e.Data.(int); ok {
			o.c += d
		}
	}
}
