package crdt

import (
	"testing"
)

func TestAddVertex(t *testing.T) {
	g := Init(NewGraph(), graphHandlers, graphQueries)
	g.Process(AddVertexCmd, "a")
	g.Process(AddVertexCmd, "b")
	exists, err := g.Query(ExistsVertexCmd, "a")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if exists != "true" {
		t.Errorf("Query(EXISTSV a): expected true, got %s", exists)
	}
	exists, err = g.Query(ExistsVertexCmd, "b")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if exists != "true" {
		t.Errorf("Query(EXISTSV b): expected true, got %s", exists)
	}
	neighbors, err := g.Query(NeighborsCmd, "a")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if neighbors != "{}" {
		t.Errorf("Query(NEIGH a): expected {}, got %q", neighbors)
	}
}

func TestAddEdge(t *testing.T) {
	g := Init(NewGraph(), graphHandlers, graphQueries)
	g.Process(AddVertexCmd, "a")
	g.Process(AddVertexCmd, "b")
	g.Process(AddVertexCmd, "c")
	g.Process(AddEdgeCmd, NewEdge("a", "b"))
	g.Process(AddEdgeCmd, NewEdge("a", "c"))
	neighbors, err := g.Query(NeighborsCmd, "a")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if neighbors != "{b, c}" && neighbors != "{c, b}" {
		t.Errorf("Query(NEIGH a): expected {b, c}, got %q", neighbors)
	}
	neighbors, err = g.Query(NeighborsCmd, "b")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if neighbors != "{}" {
		t.Errorf("Query(NEIGH b): expected {}, got %q", neighbors)
	}
}

func TestRemove(t *testing.T) {
	g := Init(NewGraph(), graphHandlers, graphQueries)
	g.Process(AddVertexCmd, "a")
	g.Process(AddVertexCmd, "b")
	g.Process(AddVertexCmd, "c")
	g.Process(AddEdgeCmd, NewEdge("a", "b"))
	g.Process(AddEdgeCmd, NewEdge("a", "c"))
	g.Process(RemoveVertexCmd, "a")
	exists, err := g.Query(ExistsVertexCmd, "a")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if exists != "true" {
		t.Errorf("Query(EXISTSV a), expected true, got false")
	}
	g.Process(RemoveEdgeCmd, NewEdge("a", "b"))
	g.Process(RemoveEdgeCmd, NewEdge("a", "c"))
	g.Process(RemoveVertexCmd, "a")
	exists, err = g.Query(ExistsVertexCmd, "a")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if exists != "false" {
		t.Errorf("Query(EXISTSV a), expected false, got true")
	}
}


