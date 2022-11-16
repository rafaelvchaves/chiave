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
	replica util.Replica
	state   *pb.DeltaSet
	delta   *pb.DeltaSet
}

func newState(r string) *pb.DeltaSet {
	return &pb.DeltaSet{
		Add: make(map[string]*pb.Dots),
		Rem: make(map[string]*pb.Dots),
		DVV: &pb.DVV{
			Dot: &pb.Dot{
				Replica: r,
				N:       0,
			},
		},
	}
}

func NewSet(replica util.Replica) *Set {
	return &Set{
		replica: replica,
		state:   newState(replica.String()),
		delta:   newState(replica.String()),
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
	return listToString(s.Value())
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
		s.delta.DVV = u // is using the same pointer a problem?
	case util.CC:
		s.state.DVV = util.Join(ctx.Dvv, u)
		s.delta.DVV = util.Join(ctx.Dvv, u)
	}
	dot := u.Dot
	addDots(s.state.Add, e, dot)
	addDots(s.delta.Add, e, dot)
	delete(s.state.Rem, e)
	delete(s.delta.Rem, e)
}

func getDots(elements map[string]*pb.Dots, e string) *pb.Dots {
	d, ok := elements[e]
	if !ok {
		return &pb.Dots{}
	}
	return d
}

func (s *Set) Remove(ctx *pb.Context, e string) {
	u := util.UpdateSingle(ctx.Dvv, s.state.DVV, s.replica.String())
	switch util.Compare(ctx.Dvv, u) {
	case util.LT:
		s.state.DVV = u
		s.delta.DVV = u
		dot := u.Dot
		delete(s.state.Add, e)
		delete(s.delta.Add, e)
		addDots(s.state.Rem, e, dot)
		addDots(s.delta.Rem, e, dot)
	case util.CC:
		fmt.Println("CC")
		// filter out all dots in client context's causal history?
		dots := getDots(s.state.Add, e)
		util.Filter(func(dot *pb.Dot) bool {
			return !util.ContainedIn(dot, ctx.Dvv)
		}, &dots.Dots)
		s.state.DVV = util.Join(ctx.Dvv, u)
		s.delta.DVV = util.Join(ctx.Dvv, u)
	default:
		fmt.Println("default")
	}
}

func (s *Set) PrepareEvent() *pb.Event {
	delta := s.delta
	s.delta = newState(s.replica.String())
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
	for e, d := range s.state.Add {
		if _, ok := ds.Rem[e]; ok {
			// for every element that this replica contains, but the other
			// replica has removed, we only keep the dots that are not contained
			// in the DVV of the other replica.
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
}

func (s *Set) Context() *pb.Context {
	return &pb.Context{
		Dvv: s.state.DVV,
	}
}

//lint:ignore U1000 Ignore unused warning: only used for debugging
func (s *Set) printState() {
	newline := "\n"
	var str string
	str += "Add: " + setToString(s.state.Add) + newline
	str += "Remove: " + setToString(s.state.Rem) + newline
	str += "DVV: " + util.String(s.state.DVV) + newline
	str += "Delta Add: " + setToString(s.delta.Add) + newline
	str += "Delta Remove: " + setToString(s.delta.Rem) + newline
	fmt.Println(str)
}

func listToString(lst []string) string {
	str := "{"
	for i, e := range lst {
		str += e
		if i < len(lst)-1 {
			str += ","
		}
	}
	str = str + "}"
	return str
}

func setToString(elements map[string]*pb.Dots) string {
	var elems []string
	for e, d := range elements {
		dots := d.Dots
		for _, dot := range dots {
			elems = append(elems, fmt.Sprintf("(%s, %q, %d)", e, dot.Replica, dot.N))
		}
	}
	return listToString(elems)
}
