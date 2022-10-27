package generator

import (
	"kvs/crdt"
	"kvs/crdt/delta"
	"kvs/crdt/op"
	"kvs/crdt/state"
	"kvs/util"
)

type Generator[F crdt.Flavor] interface {
	New(crdt.DataType, util.Replica) crdt.CRDT[F]
}

type Delta struct{}

func (Delta) New(dt crdt.DataType, r util.Replica) crdt.CRDT[crdt.Delta] {
	switch dt {
	case crdt.CType:
		return delta.NewCounter(r)
	default:
		return delta.NewSet(r)
	}
}

type Op struct{}

func (Op) New(dt crdt.DataType, r util.Replica) crdt.CRDT[crdt.Op] {
	switch dt {
	case crdt.CType:
		return op.NewCounter(r)
	default:
		return op.NewSet(r)
	}
}

type State struct{}

func (State) New(dt crdt.DataType, r util.Replica) crdt.CRDT[crdt.State] {
	switch dt {
	case crdt.CType:
		return state.NewCounter(r)
	default:
		return state.NewSet(r)
	}
}
