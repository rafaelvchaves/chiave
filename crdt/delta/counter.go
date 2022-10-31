package delta

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

type delta = struct {
	posDelta map[string]int64
	negDelta map[string]int64
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

func (c *Counter) GetEvent() *pb.Event {
	return &pb.Event{
		Source:   c.replica.String(),
		Datatype: pb.DT_Counter,
		Data: &pb.Event_DeltaCounter{
			DeltaCounter: &pb.DeltaCounter{
				Pos: c.pos.GetDelta(),
				Neg: c.neg.GetDelta(),
			},
		},
	}
}

func (s *Counter) PersistEvent(event *pb.Event) {
	dc := event.GetDeltaCounter()
	if dc == nil {
		fmt.Println("warning: nil delta counter encountered in PersistEvent")
		return
	}
	d := delta{
		posDelta: dc.Pos,
		negDelta: dc.Neg,
	}
	s.Merge(d)
}
