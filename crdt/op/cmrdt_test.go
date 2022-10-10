package crdt

import (
	"testing"
)

var graphHandlers = map[string]Handler[Graph]{
	"ADDV": AddVertexHandler{},
	"RMV":  RemoveVertexHandler{},
	"ADDE": AddEdgeHandler{},
	"RME":  RemoveEdgeHandler{},
}

var graphQueries = map[string]Query[Graph]{
	"NEIGH":   NeighborQuery{},
	"EXISTSV": ExistsVertexQuery{},
	"EXISTSE": ExistsEdgeQuery{},
}

var counterHandlers = map[string]Handler[Counter]{
	"INC": IncrementHandler{},
	"DEC": DecrementHandler{},
}

var counterQueries = map[string]Query[Counter]{
	"VALUE": ValueQuery{},
}

func TestAddVertex(t *testing.T) {
	g := Init(NewGraph(), graphHandlers, graphQueries)
	if err := g.Process("ADDV", "a"); err != nil {
		t.Fatalf("unexpected Process() error: %q", err)
	}
	if err := g.Process("ADDV", "b"); err != nil {
		t.Fatalf("unexpected Process() error: %q", err)
	}
	exists, err := g.Query("EXISTSV", "a")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if exists != "true" {
		t.Errorf("Query(EXISTSV a): expected true, got %s", exists)
	}
	exists, err = g.Query("EXISTSV", "b")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if exists != "true" {
		t.Errorf("Query(EXISTSV b): expected true, got %s", exists)
	}
	neighbors, err := g.Query("NEIGH", "a")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if neighbors != "{}" {
		t.Errorf("Query(NEIGH a): expected {}, got %q", neighbors)
	}
}

func TestAddEdge(t *testing.T) {
	g := Init(NewGraph(), graphHandlers, graphQueries)
	g.Process("ADDV", "a")
	g.Process("ADDV", "b")
	g.Process("ADDV", "c")
	g.Process("ADDE", NewEdge("a", "b"))
	g.Process("ADDE", NewEdge("a", "c"))
	neighbors, err := g.Query("NEIGH", "a")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if neighbors != "{b, c}" {
		t.Errorf("Query(NEIGH a): expected {b, c}, got %q", neighbors)
	}
	neighbors, err = g.Query("NEIGH", "b")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if neighbors != "{}" {
		t.Errorf("Query(NEIGH b): expected {}, got %q", neighbors)
	}
}

func TestIncrement(t *testing.T) {
	c := Init(NewCounter(), counterHandlers, counterQueries)
	c.Process("INC", 1)
	c.Process("INC", 1)
	got, err := c.Query("VALUE", "")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if got != "2" {
		t.Errorf("Query(VALUE): expected 2, got %s", got)
	}
}
