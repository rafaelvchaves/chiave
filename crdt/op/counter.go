package crdt

type Counter interface {
	Increment()
	Decrement()
	Value() int
}
