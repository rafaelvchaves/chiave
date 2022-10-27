package crdt_test

import (
	"kvs/crdt"
	"kvs/crdt/generator"
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

func assertEqual[T comparable](t *testing.T, name string, got, want T) {
	if got != want {
		t.Errorf("%s: got %v, want %v", name, got, want)
	}
}

func testFlavor[F crdt.Flavor](t *testing.T, g generator.Generator[F]) {
	v1 := g.New(crdt.CType, util.NewReplica("a", 1))
	v2 := g.New(crdt.CType, util.NewReplica("a", 2))
	v3 := g.New(crdt.CType, util.NewReplica("a", 3))

	c1 := v1.(crdt.Counter)
	c2 := v2.(crdt.Counter)
	c3 := v3.(crdt.Counter)

	incDecMN(c1, 5, 3)
	incDecMN(c2, 1, 4)
	incDecMN(c3, 3, 0)

	assertEqual(t, "c1 initial val", c1.Value(), 2)
	assertEqual(t, "c2 initial val", c2.Value(), -3)
	assertEqual(t, "c3 initial val", c3.Value(), 3)

	e1 := v1.GetEvent()
	e2 := v2.GetEvent()
	e3 := v3.GetEvent()

	v1.PersistEvent(e2)
	assertEqual(t, "c1 after merging c2", c1.Value(), -1)
	v1.PersistEvent(e3)

	v2.PersistEvent(e3)
	assertEqual(t, "c2 after merging c3", c2.Value(), 0)
	v2.PersistEvent(e1)

	v3.PersistEvent(e1)
	assertEqual(t, "c3 after merging c1", c3.Value(), 5)
	v3.PersistEvent(e2)

	assertEqual(t, "c1 final val", c1.Value(), 2)
	assertEqual(t, "c2 final val", c2.Value(), 2)
	assertEqual(t, "c3 final val", c3.Value(), 2)
}

func TestDeltaCounter(t *testing.T) {
	testFlavor[crdt.Delta](t, generator.Delta{})
}

func TestOpCounter(t *testing.T) {
	testFlavor[crdt.Op](t, generator.Op{})
}

func TestStateCounter(t *testing.T) {
	testFlavor[crdt.State](t, generator.State{})
}
