package crdt

import (
	"kvs/data"

	"github.com/google/uuid"
)

type Graph struct {
	vertices data.Set[data.Pair[Vertex, tag]]
	edges    data.Set[data.Pair[Edge, tag]]
}

type Edge = data.Pair[string, string]

type Vertex = string

type tag = string

type Graph1 interface {
	AddVertex(v Vertex)
	RemoveVertex(v Vertex)
	AddEdge(e Edge)
	RemoveEdge(e Edge)
	LookupEdge(e Edge) bool
	LookupVertex(v Vertex) bool
}

func (g *Graph) Init() {
	g.vertices = data.NewSet[data.Pair[Vertex, tag]]()
	g.edges = data.NewSet[data.Pair[Edge, tag]]()
}

func (g *Graph) LookupVertex(v Vertex) bool {
	ok := g.vertices.Exists(EqualsVertex(v))
	return ok
}

func (g *Graph) LookupEdge(e Edge) bool {
	ok := g.edges.Exists(EqualsEdge(e))
	v1, v2 := e.First, e.Second
	return ok && g.LookupVertex(v1) && g.LookupVertex(v2)
}

type AddVertexHandler struct{}

func (_ AddVertexHandler) Prepare(_ Graph, v any) (any, bool) {
	w := uuid.New().String()
	return data.NewPair(v.(Vertex), w), true
}

func (_ AddVertexHandler) Effect(g *Graph, p any) {
	g.vertices.Add(p.(data.Pair[Vertex, tag]))
}

func (g *Graph) PrepareAddVertex(v Vertex) (data.Pair[Vertex, tag], bool) {
	w := uuid.New().String()
	return data.NewPair(v, w), true
}

func (g *Graph) EffectAddVertex(p data.Pair[Vertex, tag]) {
	g.vertices.Add(p)
}
func (g *Graph) PrepareRemoveVertex(v Vertex) (R data.Set[data.Pair[Vertex, tag]], ok bool) {
	ok = g.LookupVertex(v) && !g.edges.Exists(
		func(p data.Pair[Edge, tag]) bool { return p.First.First == v }, // ensure no edges are coming out of v
	)
	if !ok {
		return
	}
	ok = true
	R = g.vertices.Filter(EqualsVertex(v))
	return
}

func (g *Graph) EffectRemoveVertex(R data.Set[data.Pair[Vertex, tag]]) {
	g.vertices.Subtract(R)
}

func (g *Graph) PrepareAddEdge(e Edge) (p data.Pair[Edge, tag], ok bool) {
	if !g.LookupVertex(e.First) {
		return
	}
	ok = true
	w := uuid.New().String()
	p = data.NewPair(e, w)
	return
}

func (g *Graph) EffectAddEdge(p data.Pair[Edge, tag]) {
	g.edges.Add(p)
}

func (g *Graph) PrepareRemoveEdge(e Edge) (R data.Set[data.Pair[Edge, tag]], ok bool) {
	if !g.LookupEdge(e) {
		return
	}
	ok = true
	R = g.edges.Filter(EqualsEdge(e))
	return
}

func (g *Graph) EffectRemoveEdge(R data.Set[data.Pair[Edge, tag]]) {
	g.edges.Subtract(R)
}

func EqualsVertex(v Vertex) func(data.Pair[Vertex, tag]) bool {
	return func(p data.Pair[Vertex, tag]) bool { return p.First == v }
}

func EqualsEdge(e Edge) func(data.Pair[Edge, tag]) bool {
	return func(p data.Pair[Edge, tag]) bool { return p.First == e }
}
