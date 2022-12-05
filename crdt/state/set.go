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
	for _, d := range dots {
		if containsDot(elements, e, d) {
			continue
		}
		elements[e].Dots = append(elements[e].Dots, &pb.Dot{Replica: d.Replica, N: d.N})
	}
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
		dots.Dots = util.Filter2(func(dot *pb.Dot) bool {
			return !util.ContainedIn(dot, ctx.Dvv)
		}, dots.Dots)
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
		cpy[e] = &pb.Dots{
			Dots: make([]*pb.Dot, len(dots)),
		}
		copy(cpy[e].Dots, dots)
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

func containsDot(m map[string]*pb.Dots, e string, dot *pb.Dot) bool {
	for _, d := range m[e].GetDots() {
		if d.Replica == dot.Replica && d.N == dot.N {
			return true
		}
	}
	return false
}

func (s *Set) PersistEvent(event *pb.Event) {
	ds := event.GetStateSet()
	if ds == nil {
		fmt.Println("warning: nil delta set encountered in PersistEvent")
		return
	}
	for e, d := range s.state.Add {
		if _, ok := ds.Rem[e]; ok {
			d.Dots = util.Filter2(func(dot *pb.Dot) bool {
				return !util.ContainedIn(dot, ds.DVV)
			}, d.Dots)
		}
	}
	for e, d := range ds.Add {
		if _, ok := s.state.Rem[e]; ok {
			d.Dots = util.Filter2(func(dot *pb.Dot) bool {
				return !util.ContainedIn(dot, s.state.DVV)
			}, d.Dots)
		}
		addDots(s.state.Add, e, d.Dots...)
	}
	for e, d := range ds.Rem {
		dots := d.Dots
		if _, ok := s.state.Add[e]; ok {
			dots = util.Filter2(func(dot *pb.Dot) bool {
				return !util.ContainedIn(dot, s.state.DVV)
			}, dots)
		}
		addDots(s.state.Rem, e, dots...)
	}
	for _, d := range s.state.Add {
		d.Dots = util.Filter2(func(dot1 *pb.Dot) bool {
			return !exists(d.Dots, func(dot2 *pb.Dot) bool {
				return dot1.Replica == dot2.Replica && dot1.N < dot2.N
			})
		}, d.Dots)
	}
	for _, d := range s.state.Rem {
		d.Dots = util.Filter2(func(dot1 *pb.Dot) bool {
			return !exists(d.Dots, func(dot2 *pb.Dot) bool {
				return dot1.Replica == dot2.Replica && dot1.N < dot2.N
			})
		}, d.Dots)
	}
	s.state.DVV = util.Sync(s.state.DVV, ds.DVV)
}

func exists[T comparable](lst []T, p func(T) bool) bool {
	for _, x := range lst {
		if p(x) {
			return true
		}
	}
	return false
}

func (s *Set) Context() *pb.Context {
	dvv := &pb.DVV{
		Dot: &pb.Dot{
			Replica: s.replica.String(),
			N:       0,
		},
		Clock: make(map[string]int64),
	}
	copyDVV(s.state.DVV, dvv)
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
