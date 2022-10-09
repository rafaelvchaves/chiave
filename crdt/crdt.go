package crdt

import "github.com/google/uuid"

type Event[D any] struct {
	Data D
}

type Operation[V any, D any] interface {
	Prepare(V) Event[D]
	Effect(Event[D])
}


type CmRDT struct {
	operations []Operation[any, any]
}



type Counter struct {
	i int
}

func (c *Counter) Init() { c.i = 0 }

func (c *Counter) Query(_ int) int { return c.i }

func (c *Counter) Prepare(_ int, cmd int) int { return cmd }

func (c *Counter) Effect(_ int, e Event[int]) { c.i += e.Data }

type Graph struct {
	vertices Set[Pair[Vertex, Tag]]
	edges    Set[Pair[Edge, Tag]]
}

type Pair[T1, T2 any] struct {
	fst T1
	snd T2
}

type Edge = Pair[string, string]

type Vertex = string

type Tag = string

func (g *Graph) Init() {
	g.vertices = NewSet[Pair[Vertex, Tag]]()
	g.edges = NewSet[Pair[Edge, Tag]]()
}

func (g *Graph) LookupVertex(v Vertex) bool {
	ok := g.vertices.Exists(EqualsVertex(v))
	return ok
}

func (g *Graph) LookupEdge(e Edge) bool {
	ok := g.edges.Exists(EqualsEdge(e))
	v1, v2 := e.fst, e.snd
	return ok && g.LookupVertex(v1) && g.LookupVertex(v2)
}

func (g *Graph) PrepareAddVertex(v Vertex) (Pair[Vertex, Tag], bool) {
	w := uuid.New().String()
	return Pair[Vertex, Tag]{v, w}, true
}	

func (g *Graph) EffectAddVertex(p Pair[Vertex, Tag]) {
	g.vertices.Add(p)
}
func (g *Graph) PrepareRemoveVertex(v Vertex) (R Set[Pair[Vertex, Tag]], ok bool) {
	ok = g.LookupVertex(v) && !g.edges.Exists(
		func(p Pair[Edge, Tag]) bool { return p.fst.fst == v }, // ensure no edges are coming out of v
	)
	if !ok {
		return
	}
	ok = true
	R = g.vertices.Filter(EqualsVertex(v))
	return
}

// Prepare: payload -> (operation, bool)
// Effect: operation -> unit

func (g *Graph) EffectRemoveVertex(R Set[Pair[Vertex, Tag]]) {
	g.vertices.Subtract(R)
}

func (g *Graph) PrepareAddEdge(v1, v2 Vertex) (p Pair[Edge, Tag], ok bool) {
	if !g.LookupVertex(v1) {
		return
	}
	ok = true
	e := Edge{v1, v2}
	w := uuid.New().String()
	p = Pair[Edge, Tag]{e, w}
	return
}

func (g *Graph) EffectAddEdge(p Pair[Edge, Tag]) {
	g.edges.Add(p)
}

func (g *Graph) PrepareRemoveEdge(v1, v2 Vertex) (R Set[Pair[Edge, Tag]], ok bool) {
	e := Edge{v1, v2}
	if !g.LookupEdge(e) {
		return
	}
	ok = true
	R = g.edges.Filter(EqualsEdge(e))
	return
}

func (g *Graph) EffectRemoveEdge(R Set[Pair[Edge, Tag]]) {
	g.edges.Subtract(R)
}

func EqualsVertex(v Vertex) func(Pair[Vertex, Tag]) bool {
	return func(p Pair[Vertex, Tag]) bool { return p.fst == v }
}

func EqualsEdge(e Edge) func(Pair[Edge, Tag]) bool {
	return func(p Pair[Edge, Tag]) bool { return p.fst == e }
}

