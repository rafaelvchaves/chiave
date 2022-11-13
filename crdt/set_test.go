package crdt_test

import (
	"kvs/crdt"
	"kvs/crdt/generator"
	pb "kvs/proto"
	"kvs/util"
	"sort"
	"testing"
)

func addRemove(s crdt.Set, add []string, rem []string) {
	for _, a := range add {
		s.Add(a)
	}
	for _, r := range rem {
		s.Remove(r)
	}
}

func setCompare(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	sort.Slice(s1, func(i, j int) bool { return s1[i] <= s1[j] })
	sort.Slice(s2, func(i, j int) bool { return s2[i] <= s2[j] })
	for i, x := range s1 {
		if s2[i] != x {
			return false
		}
	}
	return true
}

func testSet[F crdt.Flavor](t *testing.T, g generator.Generator[F]) {
	v1 := g.New(pb.DT_Set, util.NewReplica("a", 1))
	v2 := g.New(pb.DT_Set, util.NewReplica("a", 2))
	v3 := g.New(pb.DT_Set, util.NewReplica("a", 3))

	s1 := v1.(crdt.Set)
	s2 := v2.(crdt.Set)
	s3 := v3.(crdt.Set)

	addRemove(s1, []string{"a", "b"}, nil)
	addRemove(s2, []string{"a"}, []string{"a"})
	addRemove(s3, []string{"c", "d"}, []string{"d"})

	assertEqual(t, "c1 initial val", s1.Value(), []string{"a", "b"}, setCompare)
	assertEqual(t, "c2 initial val", s2.Value(), nil, setCompare)
	assertEqual(t, "c3 initial val", s3.Value(), []string{"c"}, setCompare)

	e1 := v1.GetEvent()
	e2 := v2.GetEvent()
	e3 := v3.GetEvent()

	v1.PersistEvent(e2)
	assertEqual(t, "s1 after merging s2", s1.Value(), []string{"a", "b"}, setCompare)
	v1.PersistEvent(e3)

	v2.PersistEvent(e3)
	assertEqual(t, "s2 after merging s3", s2.Value(), []string{"c"}, setCompare)
	v2.PersistEvent(e1)

	v3.PersistEvent(e1)
	assertEqual(t, "s3 after merging s1", s3.Value(), []string{"a", "b", "c"}, setCompare)
	v3.PersistEvent(e2)

	assertEqual(t, "s1 final val", s1.Value(), []string{"a", "b", "c"}, setCompare)
	assertEqual(t, "s2 final val", s2.Value(), []string{"a", "b", "c"}, setCompare)
	assertEqual(t, "s3 final val", s3.Value(), []string{"a", "b", "c"}, setCompare)
}

// func TestDeltaSet(t *testing.T) {
// 	testSet[crdt.Delta](t, generator.Delta{})
// }

func TestOpSet(t *testing.T) {
	testSet[crdt.Op](t, generator.Op{})
}

// func TestStateSet(t *testing.T) {
// 	testSet[crdt.State](t, generator.State{})
// }
