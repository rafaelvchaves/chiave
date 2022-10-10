package crdt

import "fmt"

type Counter int

func NewCounter() Counter {
	return Counter(0)
}

type IncrementHandler struct{}

func (IncrementHandler) Prepare(i Counter, _ any) (any, bool) {
	return 1, true
}

func (IncrementHandler) Effect(i *Counter, val any) {
	*i = *i + Counter(val.(int))
}

type DecrementHandler struct{}

func (DecrementHandler) Prepare(i Counter, _ any) (any, bool) {
	return 1, true
}

func (DecrementHandler) Effect(i *Counter, val any) {
	*i = *i - Counter(val.(int))
}

type ValueQuery struct{}

func (ValueQuery) Query(i Counter, _ any) string {
	return fmt.Sprintf("%d", i)
}
