package crdt

import (
	"kvs/util"
)

type Generator[F Flavor] interface {
	New(DataType, util.Replica) CRDT[F]
}

type DataType int

const (
	CType DataType = iota
	SType
	GType
)

type Flavor any

type Event[F Flavor] struct {
	Source util.Replica
	Type   DataType
	Key    string
	Data   any
}

type CRDT[F Flavor] interface {
	GetEvent() Event[F]
	PersistEvent(Event[F])
}

// Counters
type Counter interface {
	Value() int
	Increment()
	Decrement()
}

// Sets
type Set interface {
	Lookup(string) bool
	Add(string)
	Remove(string)
}

// Graphs
type Vertex string

type Edge struct {
	src, dest string
}

type Graph interface {
	AddVertex(v Vertex)
	RemoveVertex(v Vertex)
	AddEdge(e Edge)
	RemoveEdge(e Edge)
	LookupEdge(e Edge) bool
	LookupVertex(v Vertex) bool
}
