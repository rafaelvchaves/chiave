package crdt

import (
	"kvs/data"
)

type metadata struct {
	timestamp int
	replica   string
}

type SORSet struct {
	id     string
	set    map[string]data.Set[metadata]
	vclock GCounter
}

func NewStateORSet(id string) SORSet {
	return SORSet{
		id:     id,
		set:    make(map[string]data.Set[metadata]),
		vclock: NewGCounter(id),
	}
}

func (s *SORSet) Lookup(e string) bool {
	_, ok := s.set[e]
	return ok
}

func (s *SORSet) Add(e string) {
	c := s.vclock.Value() + 1
	s.vclock.Increment()
	md := metadata{
		timestamp: c,
		replica:   s.id,
	}
	mds, ok := s.set[e]
	if !ok {
		s.set[e] = data.NewSet(md)
		return
	}
	mds.RemoveWhere(func(md metadata) bool {
		return md.timestamp < c // remove all entries older than new timestamp
	})
	mds.Add(md)
}

func (s *SORSet) Remove(e string) {
	delete(s.set, e)
}

func union(a, b map[string]data.Set[metadata]) map[string]data.Set[metadata] {
	m := make(map[string]data.Set[metadata])
	for e, mdsa := range a {
		mdsb, ok := b[e]
		if !ok {
			m[e] = mdsa
			continue
		}
		m[e] = data.Union(mdsa, mdsb)
	}
	for e, mdsb := range b {
		if _, ok := m[e]; !ok {
			m[e] = mdsb
		}
	}
	return m
}

func addSafe(m map[string]data.Set[metadata], e string, md metadata) {
	old, ok := m[e]
	if !ok {
		m[e] = data.NewSet(md)
		return
	}
	old.Add(md)
}

func (s *SORSet) Merge(o SORSet) {
	U := make(map[string]data.Set[metadata])
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
		mu.ForEach(func(m1 metadata) {
			if mu.Exists(func(m2 metadata) bool {
				return m1.replica == m2.replica && m2.timestamp > m1.timestamp
			}) {
				// m1 is outdated: remove
				mu.Delete(m1)
			}
		})
	}
	s.set = U
	s.vclock.Merge(o.vclock)
}
