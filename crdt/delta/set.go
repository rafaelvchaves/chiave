package delta

import (
	"fmt"
	"kvs/crdt"
	pb "kvs/proto"
	"kvs/util"
	"strings"
)

var _ crdt.Set = &Set{}
var _ crdt.CRDT[crdt.Delta] = &Set{}

type Set struct {
	replica  util.Replica
	elements map[string]*pb.Dots
	history  map[string]*pb.Versions
	delta    *pb.DeltaSet
}

func NewSet(replica util.Replica) *Set {
	return &Set{
		replica:  replica,
		elements: make(map[string]*pb.Dots),
		history:  make(map[string]*pb.Versions),
		delta:    newDeltaSet(replica),
	}
}

func newDeltaSet(replica util.Replica) *pb.DeltaSet {
	return &pb.DeltaSet{
		Elements: make(map[string]*pb.Dots),
		History:  make(map[string]*pb.Versions),
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
	if _, ok := s.history[r]; !ok {
		s.history[r] = &pb.Versions{}
	}
	s.history[r].Versions = append(s.history[r].Versions, c)
	s.elements[e] = &pb.Dots{
		Dots: []*pb.Dot{dot},
	}
	if _, ok := s.delta.Elements[e]; !ok {
		s.delta.Elements[e] = &pb.Dots{}
	}
	s.delta.Elements[e].Dots = append(s.delta.Elements[e].Dots, dot)
}

func printList(lst []string) {
	str := "{"
	for i, e := range lst {
		str += e
		if i < len(lst)-1 {
			str += ","
		}
	}
	str = str + "}"
	fmt.Println(str)
}

func printElements(elements map[string]*pb.Dots) {
	var elems []string
	for e, d := range elements {
		dots := d.Dots
		for _, dot := range dots {
			elems = append(elems, fmt.Sprintf("(%s, %q, %d)", e, strings.Split(dot.Replica, ",")[1], dot.N))
		}
	}
	printList(elems)
}

func printHistory(history map[string]*pb.Versions) {
	var elems []string
	for r, v := range history {
		for _, i := range v.Versions {
			elems = append(elems, fmt.Sprintf("(%q, %d)", strings.Split(r, ",")[1], i))
		}
	}
	printList(elems)
}

//lint:ignore U1000 Ignore unused warning: only used for debugging
func (s *Set) printState() {
	fmt.Println("Elements:")
	printElements(s.elements)
	fmt.Println("History:")
	printHistory(s.history)
	fmt.Println("Delta Elements:")
	printElements(s.delta.Elements)
	fmt.Println("Delta History:")
	printHistory(s.delta.History)
}

func (s *Set) Remove(ctx *pb.Context, e string) {
	if _, ok := s.elements[e]; !ok {
		s.elements[e] = &pb.Dots{}
	}
	var removedDots []*pb.Dot
	for _, d := range s.elements[e].Dots {
		if _, ok := s.delta.History[d.Replica]; !ok {
			s.delta.History[d.Replica] = &pb.Versions{}
		}
		if util.Contains(d.N, s.history[d.Replica].Versions) {
			s.delta.History[d.Replica].Versions = append(s.delta.History[d.Replica].Versions, d.N)
			removedDots = append(removedDots, d)
		}
	}
	delete(s.elements, e)
	for _, d := range removedDots {
		r, c := d.Replica, d.N
		if v, ok := s.history[r]; ok {
			versions := v.Versions
			util.Filter(func(i int64) bool { return i != c }, &versions)
			v.Versions = versions
		}
	}
}

func (s *Set) GetEvent() *pb.Event {
	delta := s.delta
	s.delta = newDeltaSet(s.replica)
	return &pb.Event{
		Source:   s.replica.String(),
		Datatype: pb.DT_Set,
		Data: &pb.Event_DeltaSet{
			DeltaSet: delta,
		},
	}
}

func (s *Set) PersistEvent(event *pb.Event) {
	ds := event.GetDeltaSet()
	if ds == nil {
		fmt.Println("warning: nil delta set encountered in PersistEvent")
		return
	}
	for e, d := range s.elements {
		dots := d.Dots
		_, ok := ds.Elements[e]
		if !ok {
			util.Filter(func(dot *pb.Dot) bool {
				v, ok := ds.History[dot.Replica]
				if !ok {
					return true
				}
				return !util.Contains(dot.N, v.Versions)
			}, &dots)
		}
	}
	for eo, do := range ds.Elements {
		dots := do.Dots
		d, ok := s.elements[eo]
		if !ok {
			util.Filter(func(dot *pb.Dot) bool {
				v, ok := ds.History[dot.Replica]
				if !ok {
					return true
				}
				return !util.Contains(dot.N, v.Versions)
			}, &dots)
			s.elements[eo] = &pb.Dots{Dots: dots}
			continue
		}
		d.Dots = append(d.Dots, dots...)
	}
	for r, v := range ds.History {
		versions := v.Versions
		if _, ok := s.history[r]; !ok {
			s.history[r] = &pb.Versions{}
		}
		s.history[r].Versions = append(s.history[r].Versions, versions...)
	}
}
