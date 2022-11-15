package crdt_test

import (
	"kvs/crdt"
	"kvs/crdt/generator"
	pb "kvs/proto"
	"kvs/util"
	"sort"
	"testing"
)

type replica[F crdt.Flavor] struct {
	id      util.Replica
	context *pb.Context
	set     crdt.Set
	data    crdt.CRDT[F]
}

func new[F crdt.Flavor](id util.Replica, g generator.Generator[F]) *replica[F] {
	data := g.New(pb.DT_Set, id)
	return &replica[F]{
		id: id,
		context: &pb.Context{},
		set:  data.(crdt.Set),
		data: data,
	}
}

func (r *replica[_]) addRemove(add []string, rem []string) {
	n := int64(0)
	for _, e := range add {
		n++
		r.set.Add(r.context, e)
		r.context.Dvv = util.Sync(r.context.Dvv, r.data.Context().Dvv)
	}
	for _, e := range rem {
		n++
		r.set.Remove(r.context, e)
		r.context.Dvv = util.Sync(r.context.Dvv, r.data.Context().Dvv)
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
	r1 := new(util.NewReplica("a", 1), g)
	r2 := new(util.NewReplica("a", 2), g)
	r3 := new(util.NewReplica("a", 3), g)

	r1.addRemove([]string{"x", "y"}, nil)
	r2.addRemove([]string{"x"}, []string{"x"})
	r3.addRemove([]string{"z", "w"}, []string{"w"})

	assertEqual(t, "c1 initial val", r1.set.Value(), []string{"x", "y"}, setCompare)
	assertEqual(t, "c2 initial val", r2.set.Value(), nil, setCompare)
	assertEqual(t, "c3 initial val", r3.set.Value(), []string{"z"}, setCompare)

	e1 := r1.data.PrepareEvent()
	e2 := r2.data.PrepareEvent()
	e3 := r3.data.PrepareEvent()

	r1.data.PersistEvent(e2)
	assertEqual(t, "s1 after merging s2", r1.set.Value(), []string{"x", "y"}, setCompare)
	r1.data.PersistEvent(e3)

	r2.data.PersistEvent(e3)
	assertEqual(t, "s2 after merging s3", r2.set.Value(), []string{"z"}, setCompare)
	r2.data.PersistEvent(e1)

	r3.data.PersistEvent(e1)
	assertEqual(t, "s3 after merging s1", r3.set.Value(), []string{"x", "y", "z"}, setCompare)
	r3.data.PersistEvent(e2)

	assertEqual(t, "s1 final val", r1.set.Value(), []string{"x", "y", "z"}, setCompare)
	assertEqual(t, "s2 final val", r2.set.Value(), []string{"x", "y", "z"}, setCompare)
	assertEqual(t, "s3 final val", r3.set.Value(), []string{"x", "y", "z"}, setCompare)
}

func TestDeltaSet(t *testing.T) {
	testSet[crdt.Delta](t, generator.Delta{})
}

func TestOpSet(t *testing.T) {
	testSet[crdt.Op](t, generator.Op{})
}

// func TestStateSet(t *testing.T) {
// 	testSet[crdt.State](t, generator.State{})
// }
