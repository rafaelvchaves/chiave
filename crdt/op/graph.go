package op

// import (
// 	"kvs/crdt"
// 	"kvs/data"

// 	"github.com/google/uuid"
// )

// type OGraph struct {
// 	vertices data.Set[taggedVertex]
// 	edges    data.Set[taggedEdge]
// }

// type Vertex = crdt.Vertex

// type Edge = crdt.Edge

type tag = string

// type taggedVertex = data.Pair[Vertex, tag]

// type taggedEdge = data.Pair[Edge, tag]

// func EqualsVertex(v Vertex) func(taggedVertex) bool {
// 	return func(p taggedVertex) bool { return p.First == v }
// }

// func EqualsEdge(e Edge) func(taggedEdge) bool {
// 	return func(p taggedEdge) bool { return p.First == e }
// }

// func NewGraph() *OGraph {
// 	return &OGraph{
// 		vertices: data.NewSet[taggedVertex](),
// 		edges:    data.NewSet[taggedEdge](),
// 	}
// }

// func (g *OGraph) LookupVertex(v Vertex) bool {
// 	ok := g.vertices.Exists(EqualsVertex(v))
// 	return ok
// }

// func (g *OGraph) LookupEdge(e Edge) bool {
// 	ok := g.edges.Exists(EqualsEdge(e))
// 	v1, v2 := e.First, e.Second
// 	return ok && g.LookupVertex(v1) && g.LookupVertex(v2)
// }

// type AddVertexHandler struct{}

// func (AddVertexHandler) Prepare(_ Graph, val any) (any, bool) {
// 	w := uuid.New().String()
// 	return data.NewPair(val.(Vertex), w), true
// }

// func (AddVertexHandler) Effect(g *Graph, p any) {
// 	g.vertices.Add(p.(data.Pair[Vertex, tag]))
// }

// type RemoveVertexHandler struct{}

// func (RemoveVertexHandler) Prepare(g Graph, val any) (R any, ok bool) {
// 	v := val.(Vertex)
// 	ok = g.LookupVertex(v) && !g.edges.Exists(
// 		func(p data.Pair[Edge, tag]) bool { return p.First.First == v }, // ensure no edges are coming out of v
// 	)
// 	if !ok {
// 		return
// 	}
// 	ok = true
// 	R = g.vertices.Filter(EqualsVertex(v))
// 	return
// }

// func (RemoveVertexHandler) Effect(g *Graph, R any) {
// 	g.vertices.Subtract(R.(data.Set[taggedVertex]))
// }

// type AddEdgeHandler struct{}

// func (AddEdgeHandler) Prepare(g Graph, val any) (p any, ok bool) {
// 	e := val.(Edge)
// 	if !g.LookupVertex(e.First) {
// 		return
// 	}
// 	ok = true
// 	w := uuid.New().String()
// 	p = data.NewPair(e, w)
// 	return
// }

// func (AddEdgeHandler) Effect(g *Graph, p any) {
// 	g.edges.Add(p.(taggedEdge))
// }

// type RemoveEdgeHandler struct{}

// func (RemoveEdgeHandler) Prepare(g Graph, val any) (R any, ok bool) {
// 	e := val.(Edge)
// 	if !g.LookupEdge(e) {
// 		return
// 	}
// 	ok = true
// 	R = g.edges.Filter(EqualsEdge(e))
// 	return
// }

// func (RemoveEdgeHandler) Effect(g *Graph, R any) {
// 	g.edges.Subtract(R.(data.Set[taggedEdge]))
// }

// type NeighborQuery struct{}

// type str string

// func (s str) String() string {
// 	return string(s)
// }

// func (NeighborQuery) Query(g Graph, args any) string {
// 	v, ok := args.(Vertex)
// 	if !ok {
// 		return "{}"
// 	}
// 	neighbors := data.NewSet[str]()
// 	g.edges.ForEach(func(p taggedEdge) {
// 		if p.First.First == v {
// 			neighbors.Add(str(p.First.Second))
// 		}
// 	})
// 	return neighbors.String()
// }

// type ExistsVertexQuery struct{}

// func (ExistsVertexQuery) Query(g Graph, args any) string {
// 	v, ok := args.(Vertex)
// 	if !ok {
// 		return String(false)
// 	}
// 	return String(g.LookupVertex(v))
// }

// type ExistsEdgeQuery struct{}

// func (ExistsEdgeQuery) Query(g Graph, args any) string {
// 	e, ok := args.(Edge)
// 	if !ok {
// 		return String(false)
// 	}
// 	return String(g.LookupEdge(e))
// }

// func String(b bool) string {
// 	if b {
// 		return "true"
// 	}
// 	return "false"
// }
