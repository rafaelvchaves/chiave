package crdt

type Counter interface {
	CRDT
	Value() int

	Increment()

	Decrement()
}