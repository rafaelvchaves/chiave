package state

import (
	"fmt"
	pb "kvs/proto"
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
	return fmt.Sprintf("%d", c.Value())
}

func (c *Counter) Merge(o Counter) {
	fmt.Println(c.String())
	fmt.Println(o.String())
	c.pos.Merge(o.pos.vec)
	c.neg.Merge(o.neg.vec)
}

func (c *Counter) Copy() Counter {
	cpy := NewCounter(c.replica)
	cpy.pos = c.pos.Copy()
	cpy.neg = c.neg.Copy()
	return *cpy
}

func (c *Counter) GetEvent() *pb.Event {
	return &pb.Event{
		Source:   c.replica.String(),
		Datatype: pb.DT_Counter,
		Data: &pb.Event_StateCounter{
			StateCounter: &pb.StateCounter{
				Pos: c.pos.Copy().vec,
				Neg: c.neg.Copy().vec,
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
	s.pos.Merge(sc.Pos)
	s.neg.Merge(sc.Neg)
}
