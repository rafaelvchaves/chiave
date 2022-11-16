package state

import (
	"fmt"
	"kvs/crdt"
	pb "kvs/proto"
	"kvs/util"
)

var _ crdt.Counter = &Counter{}
var _ crdt.CRDT[crdt.State] = &Counter{}

type Counter struct {
	replica util.Replica
	pos     *pb.GCounter
	neg     *pb.GCounter
}

func NewCounter(replica util.Replica) *Counter {
	return &Counter{
		replica: replica,
		pos:     NewGCounter(replica.String()),
		neg:     NewGCounter(replica.String()),
	}
}

func (c *Counter) Value() int64 {
	return Value(c.pos) - Value(c.neg)
}

func (c *Counter) Increment() {
	Increment(c.pos)
}

func (c *Counter) Decrement() {
	Increment(c.neg)
}

func (c Counter) String() string {
	return fmt.Sprintf("%d", c.Value())
}

func (c *Counter) Copy() Counter {
	cpy := NewCounter(c.replica)
	cpy.pos = Copy(c.pos)
	cpy.neg = Copy(c.neg)
	return *cpy
}

func (c *Counter) PrepareEvent() *pb.Event {
	return &pb.Event{
		Source:   c.replica.String(),
		Datatype: pb.DT_Counter,
		Data: &pb.Event_StateCounter{
			StateCounter: &pb.StateCounter{
				Pos: Copy(c.pos),
				Neg: Copy(c.neg),
			},
		},
	}
}

func (s *Counter) PersistEvent(event *pb.Event) {
	sc := event.GetStateCounter()
	if sc == nil {
		fmt.Println("warning: nil state counter encountered in PersistEvent")
		return
	}
	Merge(s.pos, sc.Pos)
	Merge(s.neg, sc.Neg)
}

func (c *Counter) Context() *pb.Context {
	return nil
}
