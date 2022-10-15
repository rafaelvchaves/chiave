package delta

import (
	"kvs/crdt"
	"kvs/util"
)

type Counter struct {
}

func NewCounter(replica util.Replica) *Counter {
	return &Counter{}
}

func (c *Counter) Value() int                          { return 0 }
func (c *Counter) Increment()                          {}
func (c *Counter) Decrement()                          {}
func (c *Counter) GetEvent() crdt.Event[CRDT]          { return crdt.Event[CRDT]{} }
func (c *Counter) PersistEvent(event crdt.Event[CRDT]) {}
