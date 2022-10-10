package crdt

type StateCounter struct {
	id string
	pos GCounter
	neg GCounter
}

func NewStateCounter(id string) StateCounter {
	return StateCounter{
		id: id,
		pos: NewGCounter(id),
		neg: NewGCounter(id),
	}
}

func (s *StateCounter) Value() int {
	return s.pos.Value() - s.neg.Value()
}

func (s *StateCounter) Increment() {
	s.pos.Increment()
}

func (s *StateCounter) Decrement() {
	s.neg.Increment()
}

func (s *StateCounter) Merge(o StateCounter) {
	s.pos.Merge(o.pos)
	s.neg.Merge(o.neg)
}