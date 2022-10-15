package op

import (
	"kvs/crdt"
	"kvs/util"
)

type CRDT struct{}

type Generator struct{}

func (Generator) New(dt crdt.DataType, r util.Replica) crdt.CRDT[CRDT] {
	switch dt {
	case crdt.CType:
		return NewCounter(r)
	case crdt.SType:
		return NewSet(r)
	default:
		return NewGraph(r)
	}
}
