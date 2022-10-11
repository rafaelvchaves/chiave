package crdt

type Counter interface {
	Value() int

	Increment()

	Decrement()
}