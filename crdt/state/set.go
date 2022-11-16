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
	replica util.Replica
	state   *pb.StateSet
}

func newState(r string) *pb.StateSet {
	return &pb.StateSet{
		Add: make(map[string]*pb.Dots),
		Rem: make(map[string]*pb.Dots),
		DVV: &pb.DVV{
			Dot: &pb.Dot{
				Replica: r,
				N:       0,
			},
			Clock: make(map[string]int64),
		},
	}
}

func NewSet(replica util.Replica) *Set {
	return &Set{
		replica: replica,
		state:   newState(replica.String()),
	}
}

func (s *Set) Value() []string {
	state := s.state
	var result []string
	for e, d := range state.Add {
		if len(d.Dots) == 0 {
			continue
		}
		result = append(result, e)
	}
	return result
}

func (s *Set) String() string {
	return util.ListToString(s.Value())
}

func addDots(elements map[string]*pb.Dots, e string, dots ...*pb.Dot) {
	if _, ok := elements[e]; !ok {
		elements[e] = &pb.Dots{}
	}
	elements[e].Dots = append(elements[e].Dots, dots...)
}

func (s *Set) Add(ctx *pb.Context, e string) {
	u := util.UpdateSingle(ctx.Dvv, s.state.DVV, s.replica.String())
	switch util.Compare(ctx.Dvv, u) {
	case util.LT:
		s.state.DVV = u
	case util.CC:
		s.state.DVV = util.Join(ctx.Dvv, u)
	}
	dot := u.Dot
	addDots(s.state.Add, e, dot)
	delete(s.state.Rem, e)
}

func (s *Set) Remove(ctx *pb.Context, e string) {
	u := util.UpdateSingle(ctx.Dvv, s.state.DVV, s.replica.String())
	switch util.Compare(ctx.Dvv, u) {
	case util.LT:
		s.state.DVV = u
		dot := u.Dot
		delete(s.state.Add, e)
		addDots(s.state.Rem, e, dot)
	case util.CC:
		dots := getDots(s.state.Add, e)
		util.Filter(func(dot *pb.Dot) bool {
			return !util.ContainedIn(dot, ctx.Dvv)
		}, &dots.Dots)
		s.state.DVV = util.Join(ctx.Dvv, u)
	}
}

func getDots(elements map[string]*pb.Dots, e string) *pb.Dots {
	d, ok := elements[e]
	if !ok {
		return &pb.Dots{}
	}
	return d
}

func copyDotMap(m, cpy map[string]*pb.Dots) {
	for e, d := range m {
		dots := d.Dots
		cpy[e] = &pb.Dots{}
		for _, d := range dots {
			cpy[e].Dots = append(cpy[e].Dots, &pb.Dot{
				Replica: d.Replica,
				N:       d.N,
			})
		}
	}
}

func copyDVV(dvv, cpy *pb.DVV) {
	cpy.Dot.Replica = dvv.Dot.Replica
	cpy.Dot.N = dvv.Dot.N
	for r, c := range dvv.Clock {
		cpy.Clock[r] = c
	}
}

func (s *Set) copy() *pb.StateSet {
	cpy := newState(s.replica.String())
	copyDotMap(s.state.Add, cpy.Add)
	copyDotMap(s.state.Rem, cpy.Rem)
	copyDVV(s.state.DVV, cpy.DVV)
	return cpy
}

func (s *Set) PrepareEvent() *pb.Event {
	return &pb.Event{
		Source:   s.replica.String(),
		Datatype: pb.DT_Set,
		Data: &pb.Event_StateSet{
			StateSet: s.copy(),
		},
	}
}

func (s *Set) PersistEvent(event *pb.Event) {
	ds := event.GetStateSet()
	if ds == nil {
		fmt.Println("warning: nil delta set encountered in PersistEvent")
		return
	}
	for e, d := range s.state.Add {
		if _, ok := ds.Rem[e]; ok {
			util.Filter(func(dot *pb.Dot) bool {
				return !util.ContainedIn(dot, ds.DVV)
			}, &d.Dots)
		}
	}
	for e, d := range ds.Add {
		dots := d.Dots
		if _, ok := s.state.Rem[e]; ok {
			util.Filter(func(dot *pb.Dot) bool {
				return !util.ContainedIn(dot, s.state.DVV)
			}, &dots)
		}
		addDots(s.state.Add, e, dots...)
	}
	for e, d := range ds.Rem {
		dots := d.Dots
		if _, ok := s.state.Add[e]; ok {
			util.Filter(func(dot *pb.Dot) bool {
				return !util.ContainedIn(dot, s.state.DVV)
			}, &dots)
		}
		addDots(s.state.Rem, e, dots...)
	}
	s.state.DVV = util.Sync(s.state.DVV, ds.DVV)
	// s.printState(fmt.Sprintf("worker %s after persistevent from %s", s.replica, event.Source))
}

func (s *Set) Context() *pb.Context {
	return &pb.Context{
		Dvv: s.state.DVV,
	}
}

//lint:ignore U1000 Ignore unused warning: only used for debugging
func (s *Set) printState(header string) {
	newline := "\n"
	str := header + newline
	str += "Add: " + setToString(s.state.Add) + newline
	str += "Remove: " + setToString(s.state.Rem) + newline
	str += "DVV: " + util.String(s.state.DVV) + newline
	fmt.Println(str)
}

func setToString(elements map[string]*pb.Dots) string {
	var elems []string
	for e, d := range elements {
		dots := d.Dots
		for _, dot := range dots {
			elems = append(elems, fmt.Sprintf("(%s, %q, %d)", e, dot.Replica, dot.N))
		}
	}
	return util.ListToString(elems)
}

// type Set struct {
// 	replica  util.Replica
// 	elements map[string]*pb.Dots
// 	history  map[string]int64
// }

// func NewSet(replica util.Replica) *Set {
// 	return &Set{
// 		replica:  replica,
// 		elements: make(map[string]*pb.Dots),
// 		history:  make(map[string]int64),
// 	}
// }

// func (s *Set) Value() []string {
// 	var result []string
// 	for e, d := range s.elements {
// 		if len(d.Dots) == 0 {
// 			continue
// 		}
// 		result = append(result, e)
// 	}
// 	return result
// }

// func (s *Set) String() string {
// 	set := s.Value()
// 	str := "{"
// 	for i, e := range set {
// 		str += e
// 		if i < len(s.elements) {
// 			str += ","
// 		}
// 	}
// 	return str + "}"
// }

// func (s *Set) Add(ctx *pb.Context, e string) {
// 	dot := ctx.Dvv.Dot
// 	r, c := dot.Replica, dot.N
// 	s.history[r] = c
// 	s.elements[e] = &pb.Dots{
// 		Dots: []*pb.Dot{dot},
// 	}
// }

// func (s *Set) Remove(ctx *pb.Context, e string) {
// 	delete(s.elements, e)
// }

// func (s *Set) PrepareEvent() *pb.Event {
// 	return &pb.Event{
// 		Source:   s.replica.String(),
// 		Datatype: pb.DT_Set,
// 		Data: &pb.Event_StateSet{
// 			StateSet: &pb.StateSet{
// 				Elements: s.elements,
// 				History:  s.history,
// 			},
// 		},
// 	}
// }

// func (s *Set) PersistEvent(event *pb.Event) {
// 	ss := event.GetStateSet()
// 	if ss == nil {
// 		fmt.Println("warning: nil state set encountered in PersistEvent")
// 		return
// 	}
// 	for e, d := range s.elements {
// 		dots := d.Dots
// 		_, ok := ss.Elements[e]
// 		if !ok {
// 			util.Filter(func(dot *pb.Dot) bool {
// 				return dot.N > ss.History[dot.Replica]
// 			}, &dots)
// 		}
// 	}
// 	for eo, do := range ss.Elements {
// 		dots := do.Dots
// 		d, ok := s.elements[eo]
// 		if !ok {
// 			util.Filter(func(dot *pb.Dot) bool {
// 				return dot.N > s.history[dot.Replica]
// 			}, &dots)
// 			s.elements[eo] = &pb.Dots{Dots: dots}
// 			continue
// 		}
// 		d.Dots = append(d.Dots, dots...)
// 	}
// 	for e, d := range s.elements {
// 		if len(d.Dots) == 0 {
// 			delete(s.elements, e)
// 			continue
// 		}
// 		maxDot := &pb.Dot{}
// 		for _, dot := range d.Dots {
// 			if dot.N > maxDot.N {
// 				maxDot = dot
// 			}
// 		}
// 		s.elements[e] = &pb.Dots{Dots: []*pb.Dot{maxDot}}
// 	}
// 	for e, vo := range ss.History {
// 		v := safeGet(s.history, e)
// 		s.history[e] = util.Max(v, vo)
// 	}
// }

// func safeGet(m map[string]int64, r string) int64 {
// 	v, ok := m[r]
// 	if !ok {
// 		return 0
// 	}
// 	return v
// }

// func (s *Set) Context() *pb.Context {
// 	return nil
// }
