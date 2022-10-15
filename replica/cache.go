package replica

import "kvs/crdt"

type Cache[F crdt.Flavor] struct {
	cache map[string]crdt.CRDT[F]
}

func NewCache[F crdt.Flavor]() Store[F] {
	return &Cache[F]{
		cache: make(map[string]crdt.CRDT[F]),
	}
}

func (c *Cache[F]) Get(key string) (crdt.CRDT[F], bool) {
	v, ok := c.cache[key]
	return v, ok
}

func (c *Cache[F]) GetOrDefault(key string, def crdt.CRDT[F]) crdt.CRDT[F] {
	v, ok := c.cache[key]
	if !ok {
		return def
	}
	return v
}

func (c *Cache[F]) Put(key string, value crdt.CRDT[F]) {
	c.cache[key] = value
}
