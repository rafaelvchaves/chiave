package replica

import "kvs/crdt"

type Store interface {
	Get(string) (crdt.CRDT, bool)
	Put(string, crdt.CRDT)
}
