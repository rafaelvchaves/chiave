package delta

import (
	"kvs/crdt"
	"kvs/crdt/state"
	pb "kvs/proto"
	"kvs/util"
)

var _ crdt.Set = &Set{}
var _ crdt.CRDT[crdt.Delta] = &Set{}

type Set struct {
	replica util.Replica
	state   *state.Set
	delta   *state.Set
}

func NewSet(replica util.Replica) *Set {
	return &Set{
		replica: replica,
		state:   state.NewSet(replica),
		delta:   state.NewSet(replica),
	}
}

func (s *Set) Value() []string {
	return s.state.Value()
}

func (s *Set) String() string {
	return util.ListToString(s.Value())
}

func (s *Set) Add(ctx *pb.Context, e string) {
	s.delta.Add(ctx, e)
	s.state.Add(ctx, e)
}

func (s *Set) Remove(ctx *pb.Context, e string) {
	s.delta.Remove(ctx, e)
	s.state.Remove(ctx, e)
}

func (s *Set) PrepareEvent() *pb.Event {
	delta := s.delta
	s.delta = state.NewSet(s.replica)
	return delta.PrepareEvent()
}

func (s *Set) PersistEvent(event *pb.Event) {
	s.state.PersistEvent(event)
}

func (s *Set) Context() *pb.Context {
	return s.state.Context()
}
