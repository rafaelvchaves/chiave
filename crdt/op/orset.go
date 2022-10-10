package crdt

import (
	"kvs/data"

	"github.com/google/uuid"
)

type ORSet struct {
	s data.Set[data.Pair[string, tag]]
}

func NewORSet() ORSet {
	return ORSet{
		s: data.NewSet[data.Pair[string, tag]](),
	}
}

type AddHandler struct{}

func (AddHandler) Prepare(s ORSet, val any) (any, bool) {
	alpha := uuid.New().String()
	return data.NewPair(val.(string), alpha), true
}

func (AddHandler) Effect(s *ORSet, val any) {
	v, ok := val.(data.Pair[string, tag])
	if !ok {
		return
	}
	s.s.Add(v)
}

type RemoveHandler struct{}

func (RemoveHandler) Prepare(s ORSet, val any) (any, bool) {
	R := s.s.Filter(func(p data.Pair[string, tag]) bool {
		return p.First == val.(string)
	})
	return R, true
}

func (RemoveHandler) Effect(s *ORSet, val any) {
	v, ok := val.(data.Set[data.Pair[string, tag]])
	if !ok {
		return
	}
	s.s.Subtract(v)
}

type ExistsQuery struct{}

func (ExistsQuery) Query(s ORSet, args any) string {
	return String(s.s.Exists(func(p data.Pair[string, tag]) bool { return p.First == args.(string) }))
}
