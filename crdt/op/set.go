package op

import (
	"kvs/crdt"
	"kvs/util"
)

type Set struct {
	replica util.Replica
	s       util.Set[util.Pair[string, tag]]
}

func NewSet(r util.Replica) *Set {
	return &Set{
		replica: r,
		s:       util.NewSet[util.Pair[string, tag]](),
	}
}

func (s *Set) Add(e string)                        {}
func (s *Set) Remove(e string)                     {}
func (s *Set) Lookup(e string) bool                { return false }
func (s *Set) GetEvent() crdt.Event[CRDT]          { return crdt.Event[CRDT]{} }
func (s *Set) PersistEvent(event crdt.Event[CRDT]) {}

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
