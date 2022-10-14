package crdt

type Flavor int

const (
	Op Flavor = iota
	State
	Delta
)

type Event struct {
	Source string
	Key    string
	Data   any
}

type CRDT interface {
	GetEvent() Event
	PersistEvents([]Event)
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
