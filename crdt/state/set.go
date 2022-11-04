package state

import (
	"fmt"
	"kvs/crdt"
	pb "kvs/proto"
	"kvs/util"
)

type Set struct {
	replica util.Replica
	add     map[string]*pb.GCounter
	rem     map[string]*pb.GCounter
}

var _ crdt.Set = &Set{}

func NewSet(replica util.Replica) *Set {
	return &Set{
		replica: replica,
		add:     make(map[string]*pb.GCounter),
		rem:     make(map[string]*pb.GCounter),
	}
}

func (s *Set) Lookup(e string) bool {
	// _, ok := s.set[e]
	return true
}

func (s *Set) Value() map[string]struct{} {
	set := make(map[string]struct{})
	for k := range s.add {
		set[k] = struct{}{}
	}
	for k, vr := range s.rem {
		// if the element exists in the remove set
		// with a strictly greater timestamp,
		// remove it from the resulting set.
		if va, ok := s.add[k]; ok && Compare(va, vr) == LT {
			delete(set, k)
		}
	}
	return set
}

func (s *Set) Add(e string) {
	if va, ok := s.add[e]; ok {
		Increment(va)
		delete(s.rem, e)
		return
	}
	if vr, ok := s.rem[e]; ok {
		s.add[e] = vr
		Increment(vr)
		delete(s.rem, e)
		return
	}
	vec := NewGCounter(s.replica.String())
	Increment(vec)
	s.add[e] = vec
}

func (s *Set) Remove(e string) {
	if va, ok := s.add[e]; ok {
		s.rem[e] = va
		Increment(va)
		delete(s.add, e)
		return
	}
	if vr, ok := s.rem[e]; ok {
		s.add[e] = vr
		Increment(vr)
		delete(s.add, e)
		return
	}
	vec := NewGCounter(s.replica.String())
	Increment(vec)
	s.rem[e] = vec
}

func (s *Set) Merge(add, rem map[string]*pb.GCounter) {
	for k, v := range add {
		vo, ok := s.add[k]
		if ok {
			Merge(vo, v)
			continue
		}
		s.add[k] = v
	}
	for k, v := range rem {
		vo, ok := s.rem[k]
		if ok {
			Merge(vo, v)
			continue
		}
		s.rem[k] = v
	}
	for k, vr := range s.rem {
		va, ok := s.add[k]
		if ok && Compare(vr, va) == GT {
			delete(s.add, k)
		}
	}
	for k, va := range s.add {
		vr, ok := s.add[k]
		if ok && Compare(va, vr) == GT {
			delete(s.rem, k)
		}
	}
}

func (s *Set) GetEvent() *pb.Event {
	return &pb.Event{
		Source:   s.replica.String(),
		Datatype: pb.DT_Set,
		Data: &pb.Event_StateSet{
			StateSet: &pb.StateSet{
				Add: copy(s.add),
				Rem: copy(s.rem),
			},
		},
	}
}

func copy(m map[string]*pb.GCounter) map[string]*pb.GCounter {
	result := make(map[string]*pb.GCounter)
	for k, v := range m {
		result[k] = Copy(v)
	}
	return result
}

func (s *Set) PersistEvent(event *pb.Event) {
	fmt.Printf("%q calling persistevent, current state = %s\n", s.replica, s.String())
	ss := event.GetStateSet()
	if ss == nil {
		fmt.Println("warning: nil state set encountered in PersistEvent")
		return
	}
	s.Merge(ss.Add, ss.Rem)
	fmt.Printf("%q state now = %s\n", s.replica, s.String())
}

func (s *Set) String() string {
	set := s.Value()
	str := "{"
	i := 0
	for e := range set {
		i++
		str += e
		if i < len(set) {
			str += ","
		}
	}
	return str + "}"
}
