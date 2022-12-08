package generator

import (
	"kvs/crdt"
	"kvs/crdt/delta"
	"kvs/crdt/op"
	"kvs/crdt/state"
	pb "kvs/proto"
	"kvs/util"
	"time"
)

type Generator[F crdt.Flavor] interface {
	New(pb.DT, util.Replica) crdt.CRDT[F]
	BroadcastEpoch() time.Duration
}

type Delta struct{}

func (Delta) New(dt pb.DT, r util.Replica) crdt.CRDT[crdt.Delta] {
	switch dt {
	case pb.DT_Counter:
		return delta.NewCounter(r)
	default:
		return delta.NewSet(r)
	}
}

func (Delta) BroadcastEpoch() time.Duration {
	return 1000 * time.Millisecond
}

type Op struct{}

func (Op) New(dt pb.DT, r util.Replica) crdt.CRDT[crdt.Op] {
	switch dt {
	case pb.DT_Counter:
		return op.NewCounter(r)
	default:
		return op.NewSet(r)
	}
}

func (Op) BroadcastEpoch() time.Duration {
	return 1000 * time.Millisecond
}

type State struct{}

func (State) New(dt pb.DT, r util.Replica) crdt.CRDT[crdt.State] {
	switch dt {
	case pb.DT_Counter:
		return state.NewCounter(r)
	default:
		return state.NewSet(r)
	}
}

func (State) BroadcastEpoch() time.Duration {
	return 1000 * time.Millisecond
}
