package replica

import "kvs/crdt"

type Cache struct {
	cache map[string]crdt.CRDT
}

func NewCache() *Cache{
	return &Cache{
		cache: make(map[string]crdt.CRDT),
	}
}

func (c *Cache) Get(key string) (crdt.CRDT, bool) {
	v, ok := c.cache[key]
	return v, ok
}

func (c *Cache) Put(key string, value crdt.CRDT) {
	c.cache[key] = value
}
