package crdt

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