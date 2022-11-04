package state

import (
	"kvs/crdt"
	pb "kvs/proto"
	"kvs/util"
)

type Set struct {
	replica util.Replica
	add     map[string]GCounter
	rem     map[string]GCounter
}

var _ crdt.Set = &Set{}

func NewSet(replica util.Replica) *Set {
	return &Set{
		replica: replica,
		add:     make(map[string]GCounter),
		rem:     make(map[string]GCounter),
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
		if va, ok := s.add[k]; ok && va.Compare(vr) == LT {
			delete(set, k)
		}
	}
	return set
}

func (s *Set) Add(e string) {
	if va, ok := s.add[e]; ok {
		va.Increment()
		delete(s.rem, e)
		return
	}
	if vr, ok := s.rem[e]; ok {
		s.add[e] = vr
		vr.Increment()
		delete(s.rem, e)
		return
	}
	vec := NewGCounter(s.replica)
	vec.Increment()
	s.add[e] = vec
}

func (s *Set) Remove(e string) {
	if va, ok := s.add[e]; ok {
		s.rem[e] = va
		va.Increment()
		delete(s.add, e)
		return
	}
	if vr, ok := s.rem[e]; ok {
		s.add[e] = vr
		vr.Increment()
		delete(s.add, e)
		return
	}
	vec := NewGCounter(s.replica)
	vec.Increment()
	s.rem[e] = vec
}

func (s *Set) Merge(add, rem map[string]map[string]int64) {
	for k, v := range add {
		vo, ok := s.add[k]
		if ok {
			vo.Merge(OfMap(v))
			continue
		}
		s.add[k] = OfMap(v)
	}
	for k, v := range rem {
		vo, ok := s.rem[k]
		if ok {
			vo.Merge(OfMap(v))
			continue
		}
		s.rem[k] = OfMap(v)
	}
	for k, vr := range s.rem {
		va, ok := s.add[k]
		if ok && vr.Compare(va) == GT {
			delete(s.add, k)
		}
	}
	for k, va := range s.add {
		vr, ok := s.add[k]
		if ok && va.Compare(vr) == GT {
			delete(s.rem, k)
		}
	}
}

func (s *Set) GetEvent() *pb.Event {
	return &pb.Event{}
}

func (s *Set) PersistEvent(event *pb.Event) {}

func (s *Set) String() string { return "" }
