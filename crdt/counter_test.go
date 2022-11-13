package crdt_test

import (
	"kvs/crdt"
	"kvs/crdt/generator"
	pb "kvs/proto"
	"kvs/util"
	"testing"
)

func incDecMN(c crdt.Counter, m, n int) {
	for i := 0; i < m; i++ {
		c.Increment()
	}
	for i := 0; i < n; i++ {
		c.Decrement()
	}
}

func assertEqual[T any](t *testing.T, name string, got, want T, equals func(T, T) bool) {
	if !equals(got, want) {
		t.Errorf("%s: got %v, want %v", name, got, want)
	}
}

func intCompare(i, j int) bool { return i == j }

func testCounter[F crdt.Flavor](t *testing.T, g generator.Generator[F]) {
	v1 := g.New(pb.DT_Counter, util.NewReplica("a", 1))
	v2 := g.New(pb.DT_Counter, util.NewReplica("a", 2))
	v3 := g.New(pb.DT_Counter, util.NewReplica("a", 3))

	c1 := v1.(crdt.Counter)
	c2 := v2.(crdt.Counter)
	c3 := v3.(crdt.Counter)

	incDecMN(c1, 5, 3)
	incDecMN(c2, 1, 4)
	incDecMN(c3, 3, 0)

	assertEqual(t, "c1 initial val", c1.Value(), 2, intCompare)
	assertEqual(t, "c2 initial val", c2.Value(), -3, intCompare)
	assertEqual(t, "c3 initial val", c3.Value(), 3, intCompare)

	e1 := v1.GetEvent()
	e2 := v2.GetEvent()
	e3 := v3.GetEvent()

	v1.PersistEvent(e2)
	assertEqual(t, "c1 after merging c2", c1.Value(), -1, intCompare)
	v1.PersistEvent(e3)

	v2.PersistEvent(e3)
	assertEqual(t, "c2 after merging c3", c2.Value(), 0, intCompare)
	v2.PersistEvent(e1)

	v3.PersistEvent(e1)
	assertEqual(t, "c3 after merging c1", c3.Value(), 5, intCompare)
	v3.PersistEvent(e2)

	assertEqual(t, "c1 final val", c1.Value(), 2, intCompare)
	assertEqual(t, "c2 final val", c2.Value(), 2, intCompare)
	assertEqual(t, "c3 final val", c3.Value(), 2, intCompare)
}

func TestDeltaCounter(t *testing.T) {
	testCounter[crdt.Delta](t, generator.Delta{})
}

func TestOpCounter(t *testing.T) {
	testCounter[crdt.Op](t, generator.Op{})
}

func TestStateCounter(t *testing.T) {
	testCounter[crdt.State](t, generator.State{})
}
