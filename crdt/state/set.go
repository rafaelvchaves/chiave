package state

import (
	"kvs/crdt"
	"kvs/util"
)

type metadata struct {
	timestamp int
	replica   string
}

type Set struct {
	replica util.Replica
	set     map[string]util.Set[metadata]
	vclock  GCounter
}

func NewSet(replica util.Replica) *Set {
	return &Set{
		replica: replica,
		set:     make(map[string]util.Set[metadata]),
		vclock:  NewGCounter(replica),
	}
}

func (s *Set) Lookup(e string) bool {
	_, ok := s.set[e]
	return ok
}

func (s *Set) Add(e string) {
	c := s.vclock.Value() + 1
	s.vclock.Increment()
	md := metadata{
		timestamp: c,
		replica:   s.replica.String(),
	}
	mds, ok := s.set[e]
	if !ok {
		s.set[e] = util.NewSet(md)
		return
	}
	mds.RemoveWhere(func(md metadata) bool {
		return md.timestamp < c // remove all entries older than new timestamp
	})
	mds.Add(md)
}

func (s *Set) Remove(e string) {
	delete(s.set, e)
}

func union(a, b map[string]util.Set[metadata]) map[string]util.Set[metadata] {
	m := make(map[string]util.Set[metadata])
	for e, mdsa := range a {
		mdsb, ok := b[e]
		if !ok {
			m[e] = mdsa
			continue
		}
		m[e] = util.Union(mdsa, mdsb)
	}
	for e, mdsb := range b {
		if _, ok := m[e]; !ok {
			m[e] = mdsb
		}
	}
	return m
}

func addSafe(m map[string]util.Set[metadata], e string, md metadata) {
	old, ok := m[e]
	if !ok {
		m[e] = util.NewSet(md)
		return
	}
	old.Add(md)
}

func (s *Set) Merge(o Set) {
	U := make(map[string]util.Set[metadata])
	for e, ma := range s.set {
		ma.ForEach(func(m metadata) {
			mb, ok := o.set[e]
			if (ok && mb.Contains(m)) || m.timestamp > o.vclock.SafeGet(m.replica) {
				// means that (e, timestamp, replica) is in union(M, M')
				addSafe(U, e, m)
			}
		})
	}

	for e, mb := range o.set {
		mb.ForEach(func(m metadata) {
			ma, ok := s.set[e]
			if (!ok || !ma.Contains(m)) && m.timestamp > s.vclock.SafeGet(m.replica) {
				// means that (e, timestamp, replica) is in M''
				addSafe(U, e, m)
			}
		})
	}

	for _, mu := range U {
		mu.RemoveWhere(func(m1 metadata) bool {
			return mu.Exists(func(m2 metadata) bool {
				return m1.replica == m2.replica && m2.timestamp > m1.timestamp
			})
		})
	}
	s.set = U
	s.vclock.Merge(o.vclock)
}

func (s *Set) GetEvent() crdt.Event[CRDT] {
	return crdt.Event[CRDT]{}
}

func (s *Set) PersistEvent(event crdt.Event[CRDT]) {}
