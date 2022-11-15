package delta

import (
	"fmt"
	"kvs/crdt"
	pb "kvs/proto"
	"kvs/util"
)

var _ crdt.Set = &Set{}
var _ crdt.CRDT[crdt.Delta] = &Set{}

type Set struct {
	replica  util.Replica
	elements map[string]*pb.Dots
	dvv      *pb.DVV
	delta    *pb.DeltaSet
}

func NewSet(replica util.Replica) *Set {
	return &Set{
		replica:  replica,
		elements: make(map[string]*pb.Dots),
		dvv: &pb.DVV{
			Dot: &pb.Dot{
				Replica: replica.String(),
				N:       0,
			},
			Clock: make(map[string]int64),
		},
		delta: &pb.DeltaSet{
			Elements: make(map[string]*pb.Dots),
			Dvv: &pb.DVV{
				Dot: &pb.Dot{
					Replica: replica.String(),
					N:       0,
				},
				Clock: make(map[string]int64),
			},
		},
	}
}

func (s *Set) newDeltaSet() *pb.DeltaSet {
	return &pb.DeltaSet{
		Elements: make(map[string]*pb.Dots),
		Dvv:      s.dvv,
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

func addDots(elements map[string]*pb.Dots, e string, dots ...*pb.Dot) {
	if _, ok := elements[e]; !ok {
		elements[e] = &pb.Dots{}
	}
	elements[e].Dots = append(elements[e].Dots, dots...)
}

func (s *Set) Add(ctx *pb.Context, e string) {
	u := util.UpdateSingle(ctx.Dvv, s.dvv, s.replica.String())
	switch util.Compare(ctx.Dvv, u) {
	case util.LT:
		s.dvv = u
		s.delta.Dvv = u // is using the same pointer a problem?
	case util.CC:
		s.dvv = util.Join(ctx.Dvv, u)
		s.delta.Dvv = util.Join(ctx.Dvv, u)
	}
	dot := u.Dot
	addDots(s.elements, e, dot)
	addDots(s.delta.Elements, e, dot)
	s.printState()
}

func getDots(elements map[string]*pb.Dots, e string) []*pb.Dot {
	d, ok := elements[e]
	if !ok {
		return nil
	}
	return d.Dots
}

func (s *Set) Remove(ctx *pb.Context, e string) {
	u := util.UpdateSingle(ctx.Dvv, s.dvv, s.replica.String())
	switch util.Compare(ctx.Dvv, u) {
	case util.LT:
		s.dvv = u
		s.delta.Dvv = u
		delete(s.elements, e)
		delete(s.delta.Elements, e)
	case util.CC:
		// filter out all dots in client context's causal history?
		dots := getDots(s.elements, e)
		util.Filter(func(dot *pb.Dot) bool {
			return !util.ContainedIn(dot, ctx.Dvv)
		}, &dots)
		s.dvv = util.Join(ctx.Dvv, u)
		s.delta.Dvv = util.Join(ctx.Dvv, u)
	}
	s.printState()
}

func (s *Set) PrepareEvent() *pb.Event {
	delta := s.delta
	s.delta = s.newDeltaSet()
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
		if _, ok := ds.Elements[e]; !ok {
			// for every element that this replica contains, but the other
			// replica does not, we only keep the dots that are not contained
			// in the DVV of the other replica.
			util.Filter(func(dot *pb.Dot) bool {
				return !util.ContainedIn(dot, ds.Dvv)
			}, &dots)
		}
	}
	for e, d := range ds.Elements {
		dots := d.Dots
		if _, ok := s.elements[e]; !ok {
			util.Filter(func(dot *pb.Dot) bool {
				return !util.ContainedIn(dot, s.dvv)
			}, &dots)
		}
		addDots(s.elements, e, dots...)
	}
	s.dvv = util.Sync(s.dvv, ds.Dvv)
}

func (s *Set) Context() *pb.Context {
	return &pb.Context{
		Dvv: s.dvv,
	}
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
			elems = append(elems, fmt.Sprintf("(%s, %q, %d)", e, dot.Replica, dot.N))
		}
	}
	printList(elems)
}

func printDVV(dvv *pb.DVV) {
	fmt.Println(util.String(dvv))
}

//lint:ignore U1000 Ignore unused warning: only used for debugging
func (s *Set) printState() {
	fmt.Printf("State of %d\n", s.replica.WorkerID)
	fmt.Println("Elements:")
	printElements(s.elements)
	fmt.Println("DVV:")
	printDVV(s.dvv)
	fmt.Println("Delta Elements:")
	printElements(s.delta.Elements)
	fmt.Println("Delta DVV:")
	printDVV(s.delta.Dvv)
	fmt.Println()
}
