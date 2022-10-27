package op

import (
	"kvs/crdt"
	"kvs/util"
)

type CRDT struct{}

type Generator struct{}

func (Generator) New(dt crdt.DataType, r util.Replica) crdt.CRDT[crdt.Op] {
	switch dt {
	case crdt.CType:
		return NewCounter(r)
	default:
		return NewSet(r)
	}
}
