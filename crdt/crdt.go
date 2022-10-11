package crdt

type Event struct {
	Source string // id of source replica
	Data   any    // CvRDT: state, CmRDT: update, dCvRDT: delta
}

type CRDT interface {
	GetEvent() Event
	PersistEvent(Event)
}

// Counters
type Counter interface {
	CRDT
	Value() int

	Increment()

	Decrement()
}

// Sets
type Set interface {
	CRDT
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
	CRDT
	AddVertex(v Vertex)
	RemoveVertex(v Vertex)
	AddEdge(e Edge)
	RemoveEdge(e Edge)
	LookupEdge(e Edge) bool
	LookupVertex(v Vertex) bool
}
