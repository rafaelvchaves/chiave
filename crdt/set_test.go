package crdt_test

import (
	"kvs/crdt"
	"kvs/crdt/generator"
	pb "kvs/proto"
	"kvs/util"
	"sort"
	"testing"
)

func contextWith(replica string, seqNr int64) *pb.Context {
	return &pb.Context{
		Dot: &pb.Dot{
			Replica: replica,
			N:       seqNr,
		},
	}
}

func addRemove(replica string, s crdt.Set, add []string, rem []string) {
	n := int64(0)
	for _, a := range add {
		n++
		s.Add(contextWith(replica, n), a)
	}
	for _, r := range rem {
		n++
		s.Remove(contextWith(replica, n), r)
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
	r1 := util.NewReplica("a", 1)
	r2 := util.NewReplica("a", 2)
	r3 := util.NewReplica("a", 3)
	v1 := g.New(pb.DT_Set, r1)
	v2 := g.New(pb.DT_Set, r2)
	v3 := g.New(pb.DT_Set, r3)

	s1 := v1.(crdt.Set)
	s2 := v2.(crdt.Set)
	s3 := v3.(crdt.Set)

	addRemove(r1.String(), s1, []string{"a", "b"}, nil)
	addRemove(r2.String(), s2, []string{"a"}, []string{"a"})
	addRemove(r3.String(), s3, []string{"c", "d"}, []string{"d"})

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

func TestStateSet(t *testing.T) {
	testSet[crdt.State](t, generator.State{})
}
