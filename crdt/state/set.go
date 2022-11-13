package state

import (
	"fmt"
	"kvs/crdt"
	pb "kvs/proto"
	"kvs/util"
)

var _ crdt.Set = &Set{}
var _ crdt.CRDT[crdt.State] = &Set{}

type Set struct {
	replica  util.Replica
	elements map[string]*pb.Dots
	history  map[string]int64
}

func NewSet(replica util.Replica) *Set {
	return &Set{
		replica:  replica,
		elements: make(map[string]*pb.Dots),
		history:  make(map[string]int64),
	}
}

func (s *Set) Value() []string {
	var result []string
	for e, d := range s.elements {
		if len(d.Dots) == 0 {
			continue
		}
		result = append(result, e)
	}
	return result
}

func (s *Set) String() string {
	set := s.Value()
	str := "{"
	for i, e := range set {
		str += e
		if i < len(s.elements) {
			str += ","
		}
	}
	return str + "}"
}

func (s *Set) Add(ctx *pb.Context, e string) {
	dot := ctx.Dot
	r, c := dot.Replica, dot.N
	s.history[r] = c
	s.elements[e] = &pb.Dots{
		Dots: []*pb.Dot{dot},
	}
}

func (s *Set) Remove(ctx *pb.Context, e string) {
	delete(s.elements, e)
}

func (s *Set) GetEvent() *pb.Event {
	return &pb.Event{
		Source:   s.replica.String(),
		Datatype: pb.DT_Set,
		Data: &pb.Event_StateSet{
			StateSet: &pb.StateSet{
				Elements: s.elements,
				History:  s.history,
			},
		},
	}
}

func (s *Set) PersistEvent(event *pb.Event) {
	ss := event.GetStateSet()
	if ss == nil {
		fmt.Println("warning: nil state set encountered in PersistEvent")
		return
	}
	for e, d := range s.elements {
		dots := d.Dots
		_, ok := ss.Elements[e]
		if !ok {
			util.Filter(func(dot *pb.Dot) bool {
				return dot.N > ss.History[dot.Replica]
			}, &dots)
		}
	}
	for eo, do := range ss.Elements {
		dots := do.Dots
		d, ok := s.elements[eo]
		if !ok {
			util.Filter(func(dot *pb.Dot) bool {
				return dot.N > s.history[dot.Replica]
			}, &dots)
			s.elements[eo] = &pb.Dots{Dots: dots}
			continue
		}
		d.Dots = append(d.Dots, dots...)
	}
	for e, d := range s.elements {
		if len(d.Dots) == 0 {
			delete(s.elements, e)
			continue
		}
		maxDot := &pb.Dot{}
		for _, dot := range d.Dots {
			if dot.N > maxDot.N {
				maxDot = dot
			}
		}
		s.elements[e] = &pb.Dots{Dots: []*pb.Dot{maxDot}}
	}
	for e, vo := range ss.History {
		v := safeGet(s.history, e)
		s.history[e] = util.Max(v, vo)
	}
}

func safeGet(m map[string]int64, r string) int64 {
	v, ok := m[r]
	if !ok {
		return 0
	}
	return v
}
