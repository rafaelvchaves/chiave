package crdt

import (
	"testing"
)

func queryT(t *testing.T, g CmRDT[Graph], cmd string, arg string, want string) {
	got, err := g.Query(cmd, arg)
	if err != nil {
		t.Fatalf("unexpected Query() error: %v", err)
	}
	if got != want {
		t.Errorf("Query(%s %s): got = %q, want %q", cmd, arg, got, want)
	}
}

func TestAddVertex(t *testing.T) {
	g := Init(NewGraph(), graphHandlers, graphQueries)
	g.Process(AddVertexCmd, "a")
	g.Process(AddVertexCmd, "b")
	queryT(t, g, ExistsVertexCmd, "a", "true")
	queryT(t, g, ExistsVertexCmd, "b", "false")
	queryT(t, g, NeighborsCmd, "a", "{}")
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
	queryT(t, g, ExistsVertexCmd, "a", "true")
	g.Process(RemoveEdgeCmd, NewEdge("a", "b"))
	g.Process(RemoveEdgeCmd, NewEdge("a", "c"))
	g.Process(RemoveVertexCmd, "a")
	queryT(t, g, ExistsVertexCmd, "a", "false")
}


