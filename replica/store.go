package replica

import "kvs/crdt"

type Store interface {
	Get(string) (crdt.CRDT, bool)
	GetOrDefault(string, crdt.CRDT) crdt.CRDT
	Put(string, crdt.CRDT)
}
