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
	GetEvent() *pb.Event
	PersistEvent(*pb.Event)
}

// Counters
type Counter interface {
	Value() int
	Increment()
	Decrement()
}

// Sets
type Set interface {
	Value() []string
	Add(string)
	Remove(string)
}

// Graphs
type Vertex string

type Edge struct {
	// src, dest string
}

type Graph interface {
	AddVertex(v Vertex)
	RemoveVertex(v Vertex)
	AddEdge(e Edge)
	RemoveEdge(e Edge)
	LookupEdge(e Edge) bool
	LookupVertex(v Vertex) bool
}
