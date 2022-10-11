package crdt

type Set interface {
	Lookup(string) bool
	Add(string)
	Remove(string)
}