package delta

import (
	"kvs/crdt"
	"kvs/util"
)

type Set struct{}

func NewSet(r util.Replica) *Set {
	return &Set{}
}

func (s *Set) Add(e string)                        {}
func (s *Set) Remove(e string)                     {}
func (s *Set) Lookup(e string) bool                { return false }
func (s *Set) GetEvent() crdt.Event[CRDT]          { return crdt.Event[CRDT]{} }
func (s *Set) PersistEvent(event crdt.Event[CRDT]) {}
