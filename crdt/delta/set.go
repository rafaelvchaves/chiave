package delta

import (
	pb "kvs/proto"
	"kvs/util"
)

type Set struct{}

func NewSet(r util.Replica) *Set {
	return &Set{}
}

func (s *Set) Add(e string)                 {}
func (s *Set) Remove(e string)              {}
func (s *Set) Lookup(e string) bool         { return false }
func (s *Set) GetEvent() *pb.Event           { return &pb.Event{} }
func (s *Set) PersistEvent(event *pb.Event) {}
func (s *Set) String() string               { return "" }
