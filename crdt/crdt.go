package crdt

import (
	"kvs/util"
)

type DataType int

const (
	CType DataType = iota
	SType
	GType
)

type Delta struct{}
type State struct{}
type Op struct{}

type Flavor interface {
	Delta | State | Op
}

type Event struct {
	Source util.Replica
	Type   DataType
	Key    string
	Data   any
}

type CRDT[F Flavor] interface {
	String() string
	GetEvent() Event
	PersistEvent(Event)
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
