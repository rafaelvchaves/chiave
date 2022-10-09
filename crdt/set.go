package crdt

type Set[T comparable] struct {
	s map[T]struct{}
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{
		s: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(v T) {
	s.s[v] = struct{}{}
}

func (s *Set[T]) Delete(v T) {
	delete(s.s, v)
}

func (s *Set[T]) Contains(v T) bool {
	_, ok := s.s[v]
	return ok
}

func (s *Set[T]) Exists(f func(T) bool) bool {
	for k := range s.s {
		if f(k) {
			return true
		}
	}
	return false
}

func (s *Set[T]) Filter(f func(T) bool) Set[T] {
	filtered := NewSet[T]()
	for k := range s.s {
		if f(k) {
			filtered.Add(k)
		}
	}
	return filtered
}

func (s *Set[T]) Union(other Set[T]) {
	for k := range other.s {
		s.s[k] = struct{}{}
	}
}

func (s *Set[T]) Subtract(other Set[T]) {
	for k := range s.s {
		if _, ok := other.s[k]; ok {
			delete(s.s, k)
		}
	}
}
