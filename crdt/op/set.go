package op

import (
	pb "kvs/proto"
	"kvs/util"
	"strings"

	"github.com/emirpasic/gods/sets/treeset"
)

type Ord int

type Context interface {
	Compare(Context) Ord
	Merge(Context) Context
}

type Set struct {
	replica  util.Replica
	add, rem *treeset.Set
}

type element struct {
	e   string
	ctx Context
}

func compare(e1, e2 any) int {
	return strings.Compare(e1.(element).e, e2.(element).e)
}

func NewSet(r util.Replica) *Set {
	return &Set{
		replica: r,
		add:     treeset.NewWith(compare),
		rem:     treeset.NewWith(compare),
	}
}

func (s *Set) Add(e string, ctx Context) {
	s.add.Add(element{
		e:   e,
		ctx: ctx,
	})
}
func (s *Set) Remove(e string, ctx Context) {
	
}
func (s *Set) Lookup(e string, ctx Context) bool { return false }
func (s *Set) GetEvent() *pb.Event               { return &pb.Event{} }
func (s *Set) PersistEvent(event *pb.Event)      {}

func (s *Set) String() string { return "" }

// type AddHandler struct{}

// func (AddHandler) Prepare(s ORSet, val any) (any, bool) {
// 	alpha := uuid.New().String()
// 	return data.NewPair(val.(string), alpha), true
// }

// func (AddHandler) Effect(s *ORSet, val any) {
// 	v, ok := val.(data.Pair[string, tag])
// 	if !ok {
// 		return
// 	}
// 	s.s.Add(v)
// }

// type RemoveHandler struct{}

// func (RemoveHandler) Prepare(s ORSet, val any) (any, bool) {
// 	R := s.s.Filter(func(p data.Pair[string, tag]) bool {
// 		return p.First == val.(string)
// 	})
// 	return R, true
// }

// func (RemoveHandler) Effect(s *ORSet, val any) {
// 	v, ok := val.(data.Set[data.Pair[string, tag]])
// 	if !ok {
// 		return
// 	}
// 	s.s.Subtract(v)
// }

// type ExistsQuery struct{}

// func (ExistsQuery) Query(s ORSet, args any) string {
// 	return String(s.s.Exists(func(p data.Pair[string, tag]) bool { return p.First == args.(string) }))
// }
