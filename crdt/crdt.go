package crdt

import (
	pb "kvs/proto"
)

type Delta struct{}
type State struct{}
type Op struct{}

type Flavor interface {
	Delta | State | Op
}

type CRDT[F Flavor] interface {
	String() string
	PrepareEvent() *pb.Event
	PersistEvent(*pb.Event)
	Context() *pb.Context
}

// Counters
type Counter interface {
	Value() int64
	Increment()
	Decrement()
}

// Sets
type Set interface {
	Value() []string
	Add(*pb.Context, string)
	Remove(*pb.Context, string)
}
