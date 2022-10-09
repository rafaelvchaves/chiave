package crdt

import "github.com/google/uuid"

type Event[D any] struct {
	Data D
}

type CmRDT[V any, S any, C any, D any] interface {
	// S is the internal state: it's what you get when you query, and also is updated in Effect.
	// V??
	// C is the command type
	Init()
	Query(v V) S
	Prepare(v V, cmd C) D // side effect free
	Effect(v V, e Event[D])
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

type Opts struct {
	Kind string
}

func (g *Graph) ContainsVertex(v Vertex) bool {
	ok := g.vertices.Exists(EqVertex(v))
	return ok
}


func (g *Graph) ContainsEdge(v1, v2 Vertex) bool {
	e := Edge{v1, v2}
	ok := g.edges.Exists(EqEdge(e))
	return ok && g.ContainsVertex(v1) && g.ContainsVertex(v2)
}

func (g *Graph) PrepareAddVertex(v string) (string, bool) {
	return uuid.New().String(), true
}

func (g *Graph) EffectAddVertex(v string, w string) {
	g.vertices.Add(Pair[string, string]{v, w})
}

func EqVertex(v Vertex) func(Pair[Vertex, Tag]) bool {
	return func(p Pair[Vertex, Tag]) bool { return p.fst == v }
}

func EqEdge(e Edge) func(Pair[Edge, Tag]) bool {
	return func(p Pair[Edge, Tag]) bool { return p.fst == e }
}

func (g *Graph) PrepareRemoveVertex(v string) (s Set[Pair[Vertex, Tag]], ok bool) {
	ok = g.ContainsVertex(v) && !g.edges.Exists(
		func (p Pair[Edge, Tag]) bool { return p.fst.fst == v },
	)
	if !ok {
		return
	}
	s = g.vertices.Filter(EqVertex(v))
	return
}

func (g *Graph) EffectRemoveVertex(v string) {

}

func (g *Graph) AddEdge(v1, v2 string) {

}

func (g *Graph) RemoveEdge(v1, v2 string) {

}

func test() {
	for {
		// listen for updates

		// listen for client requests

		// send updates to others

		//
	}
}

type MyCounter struct {
	i int
}

func (m *MyCounter) Increment() {
	m.i++
}

func (m *MyCounter) Decrement() {
	m.i--
}

func (m *MyCounter) Value() int {
	return m.i
}

type CvRDT[S any, V any] interface {
	Value() V
	Merge(so S)
}

func Init(n int) {

}
