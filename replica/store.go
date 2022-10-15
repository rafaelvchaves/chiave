package replica

import "kvs/crdt"

type Store[F crdt.Flavor] interface {
	Get(string) (crdt.CRDT[F], bool)
	GetOrDefault(string, crdt.CRDT[F]) crdt.CRDT[F]
	Put(string, crdt.CRDT[F])
}
