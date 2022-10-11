package main

import (
	"fmt"
	"kvs/crdt"
	"kvs/crdt/state"
)

// Helper function for getting key from map
func Get[T crdt.CRDT](m map[string]crdt.CRDT, k string) T {
	var zero T
	value, ok := m[k]
	if !ok {
		return zero
	}
	return value.(T)
}

// Helper function for putting key in map
func Put[T crdt.CRDT](m map[string]crdt.CRDT, k string, v T) {
	m[k] = v
}

// Helper function for displaying counter value in map
func DisplayValue(name string, m map[string]crdt.CRDT, k string) {
	v := Get[crdt.Counter](m, k).Value()
	fmt.Println(fmt.Sprintf("%s value of %s: %d", name, k, v))
}

func main() {
	// Create two counters, c1 and c2.
	// c1 will locally increment 4 times, c2 will locally decrement twice.
	// we send c1's "event" to c2, and c2's "event" to c1:
	//	- if c1 and c2 are state-based, the event is their entire state
	//  - if c1 and c2 are op-based, the event is simply the net number of
	//    increments/decrements since the last event (i.e. +4 for c1 and -2 for c2).
	// Both counters call the PersistEvent() method to persist the changes locally.
	// After this call, both should return a counter value of (4 - 2) = 2.
	// c1 := op.NewCounter("c1")
	// c2 := op.NewCounter("c2")

	c1 := state.NewCounter("c1")
	c2 := state.NewCounter("c2")
	m1 := map[string]crdt.CRDT{
		"key1": c1,
	}
	m2 := map[string]crdt.CRDT{
		"key1": c2,
	}
	DisplayValue("m1", m1, "key1")
	DisplayValue("m2", m2, "key1")
	c1.Increment()
	c1.Increment()
	c1.Increment()
	c1.Increment()
	e1 := c1.GetEvent()

	c2.Decrement()
	c2.Decrement()
	e2 := c2.GetEvent()

	c1.PersistEvent(e2)
	c2.PersistEvent(e1)

	Put(m1, "key1", c1)
	Put(m2, "key1", c2)

	DisplayValue("m1", m1, "key1")
	DisplayValue("m2", m2, "key1")
}
