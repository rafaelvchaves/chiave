package data

type Element interface {
	comparable
}

type Set[T Element] struct {
	s map[T]struct{}
}

func NewSet[T Element](elements ...T) Set[T] {
	s := make(map[T]struct{})
	for _, e := range elements {
		s[e] = struct{}{}
	}
	return Set[T]{
		s: s,
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

func Union[T comparable](a, b Set[T]) Set[T] {
	result := NewSet[T]()
	for k := range a.s {
		result.s[k] = struct{}{}
	}
	for k := range b.s {
		result.s[k] = struct{}{}
	}
	return result
}

func (s *Set[T]) Subtract(other Set[T]) {
	for k := range s.s {
		if _, ok := other.s[k]; ok {
			delete(s.s, k)
		}
	}
}

func (s *Set[T]) Intersect(other Set[T]) {
	for k := range s.s {
		if _, ok := other.s[k]; !ok {
			delete(s.s, k)
		}
	}
}

func (s *Set[T]) ForEach(f func(T)) {
	for k := range s.s {
		f(k)
	}
}

func (s *Set[T]) Range(f func(T) bool) {
	for k := range s.s {
		if !f(k) {
			return
		}
	}
}

func (s *Set[T]) RemoveWhere(f func(T) bool) {
	for k := range s.s {
		if f(k) {
			delete(s.s, k)
		}
	}
}